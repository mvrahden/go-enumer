package gen

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"log"
	"regexp"
	"sort"
	"strings"

	"github.com/mvrahden/go-enumer/config"
	"golang.org/x/tools/go/packages"
)

func init() {
	reg, err := regexp.Compile("^// Code generated .* DO NOT EDIT.$")
	if err != nil {
		log.Fatalf("failed evaluating regexp. err: %s", err)
	}
	matchGeneratedFileRegex = reg
}

var matchGeneratedFileRegex *regexp.Regexp

type inspector struct {
	cfg *config.Options
}

func NewInspector(cfg *config.Options) *inspector {
	return &inspector{cfg}
}

func (i inspector) Inspect(pkg *packages.Package) (*File, error) {
	out := &File{Imports: []*Import{}}

	i.preparePackage(pkg)
	i.inspectHeader(pkg, out)
	if err := i.inspectValues(pkg, out); err != nil {
		return nil, err
	}
	if err := i.validateValues(out); err != nil {
		return nil, err
	}
	i.inspectImports(out)
	return out, nil
}

func (i inspector) inspectImports(f *File) {
	f.Imports = append(f.Imports, &Import{Path: "fmt"})
	for _, v := range i.cfg.Serializers {
		switch v {
		case "json":
			f.Imports = append(f.Imports, &Import{Path: "encoding/json"})
		case "sql":
			f.Imports = append(f.Imports, &Import{Path: "database/sql/driver"})
		}
	}
}

func (i inspector) validateValues(f *File) error {
	sort.SliceStable(f.ValueSpecs, func(i, j int) bool {
		return f.ValueSpecs[i].Value < f.ValueSpecs[j].Value
	})
	if len(f.ValueSpecs) > 0 &&
		!(f.ValueSpecs[0].Value == 0 || f.ValueSpecs[0].Value == 1) {
		return fmt.Errorf("Invalid enum set: Enums need to start with either 0 or 1.")
	}
	// ensure we have a linearly incrementing sequence of values.
	// however, an enum can assign a numeric value multiple times.
	// therefore we must only dismiss distances > 1.
	for idx := 1; idx < len(f.ValueSpecs); idx++ {
		delta := f.ValueSpecs[idx].Value - f.ValueSpecs[idx-1].Value
		if delta > 1 {
			return fmt.Errorf("Invalid enum set: Enums must be a continuous sequence with linear increments of 1.")
		}
	}
	return nil
}

// preparePackage reduce number files by dropping generated files
func (i inspector) preparePackage(pkg *packages.Package) {
	var dropIdx []int
	for idx, file := range pkg.Syntax {
		if i.isGeneratedFile(file) {
			dropIdx = append(dropIdx, idx)
		}
	}
	for adjBy, idx := range dropIdx { // drop and maintain order of files
		idx -= adjBy
		pkg.Syntax = append(pkg.Syntax[:idx], pkg.Syntax[idx+1:]...)
	}
}

func (inspector) inspectHeader(pkg *packages.Package, out *File) {
	out.Header.Package = Package{
		Name: pkg.Name,
		Path: pkg.PkgPath,
	}
	if pkg.Module == nil {
		return
	}
	out.Header.Module = Module{
		Module:       pkg.Module.Version,
		Path:         pkg.Module.Dir,
		GoModPath:    pkg.Module.GoMod,
		GoModVersion: pkg.Module.Version,
		GoVersion:    pkg.Module.GoVersion,
	}
}

func (i inspector) inspectValues(pkg *packages.Package, out *File) error {
	specs := i.determinePackageScopedValueSpecs(pkg.Syntax, out)
	for idx, s := range specs {
		vspec, err := i.evaluateValueSpec(idx, s, pkg)
		if err != nil {
			return err
		}
		out.ValueSpecs = append(out.ValueSpecs, vspec)
	}
	return nil
}

// Naive check whether file contains a "Code generated" comment.
// This function does not verify the position of the comment.
func (i inspector) isGeneratedFile(f *ast.File) bool {
	if len(f.Comments) > 0 &&
		len(f.Comments[0].List) > 0 {
		firstComment := f.Comments[0].List[0]
		if matchGeneratedFileRegex.Match([]byte(firstComment.Text)) {
			return true
		}
	}
	return false
}

type valueSpec struct {
	Type  *ast.ValueSpec
	Value *ast.Ident
}

func (i *inspector) determinePackageScopedValueSpecs(files []*ast.File, out *File) []valueSpec {
	var typeSpecs []valueSpec
	for _, f := range files {
		for _, v := range f.Decls {
			decl, ok := v.(*ast.GenDecl)
			if !ok || decl.Tok != token.CONST {
				// We only care about const declarations.
				continue
			}
			var prevType *ast.Ident // for blocks with implicit typing
			for _, spec := range decl.Specs {
				vspec, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				if vspec.Type == nil && len(vspec.Values) > 0 {
					continue
				}
				if vspec.Type != nil {
					var ok bool
					ident, ok := vspec.Type.(*ast.Ident)
					if !ok || ident == nil {
						// not the type we're searching for (as per configuration)
						continue
					}
					prevType = ident
				}
				if prevType == nil || prevType.Name != i.cfg.TypeAliasName {
					prevType = nil
					continue
				}
				if vspec.Type == nil {
					// for those with implicit type, assign the previous type
					vspec.Type = prevType
				}

				for _, v := range vspec.Names {
					if v.Name == "_" {
						// blank identifier, not what we're interested in
						continue
					}
					typeSpecs = append(typeSpecs, valueSpec{vspec, v})
				}
			}
		}
	}
	return typeSpecs
}

func (i inspector) evaluateValueSpec(idx int, v valueSpec, pkg *packages.Package) (*ValueSpec, error) {
	typ, err := i.determineTypeOfExpr(v.Type.Type)
	if err != nil {
		return nil, err
	}
	val, valStr, err := i.determineValueOfExpr(v.Value, pkg)
	if err != nil {
		return nil, err
	}
	name := v.Value.Name
	if strings.HasPrefix(v.Value.Name, i.cfg.TypeAliasName) {
		if len(i.cfg.TypeAliasName) == len(v.Value.Name) {
			return nil, fmt.Errorf("cannot determine name after trimming prefix (enum value equals type name). make sure to give a meaningful names to your enum values")
		}
		name = v.Value.Name[len(i.cfg.TypeAliasName):]
	}
	return &ValueSpec{
		Index:          idx,
		IdentifierName: v.Value.Name,
		EnumString:     name,
		Type:           typ,
		Value:          val,
		ValueString:    valStr,
	}, nil
}

func (i inspector) determineTypeOfExpr(e ast.Expr) (GoType, error) {
	switch t := e.(type) {
	case *ast.Ident:
		typ, ok := typeMap[t.Name]
		if ok {
			return typ, nil
		}
		if t.Obj == nil {
			break
		}
		// if possible detect underlying
		decl, ok := t.Obj.Decl.(*ast.TypeSpec)
		if ok {
			return i.determineTypeOfExpr(decl.Type)
		}
	}
	return GoTypeUnknown, fmt.Errorf("Invalid enum set: Enum type must be an integer-like type, found %q.", e)
}

func (i inspector) determineValueOfExpr(e ast.Expr, pkg *packages.Package) (uint64, string, error) {
	c, ok := e.(*ast.Ident)
	if !ok {
		return 0, "", fmt.Errorf("internal error: a value slipped our type evaluation (type: %+v)", e)
	}
	obj, ok := pkg.TypesInfo.Defs[c]
	if !ok {
		return 0, "", fmt.Errorf("no value for constant %q", c)
	}
	objT := obj.Type()
	if objT == nil {
		return 0, "", fmt.Errorf("definition type for constant %q is <nil>", c)
	}
	ul := objT.Underlying()
	if ul == nil {
		return 0, "", fmt.Errorf("underlying was expected to be a basic type, but was <nil>")
	}
	bul, ok := ul.(*types.Basic)
	if !ok {
		return 0, "", fmt.Errorf("underlying type for constant %q is not a basic type", c)
	}
	info := bul.Info()
	if info&types.IsInteger == 0 {
		return 0, "", fmt.Errorf("type %q is not an constant of type integer", i.cfg.TypeAliasName)
	}
	cobj, ok := obj.(*types.Const)
	if !ok {
		return 0, "", fmt.Errorf("internal error: a value slipped our type evaluation (type: %+v is not a const)", e)
	}
	value := cobj.Val()
	if value.Kind() != constant.Int {
		return 0, "", fmt.Errorf("can't happen: constant is not an integer %q", c)
	}
	i64, isInt := constant.Int64Val(value)
	u64, isUint := constant.Uint64Val(value)
	if !isInt && !isUint {
		return 0, "", fmt.Errorf("internal error: value of %s is not an integer: %s", c, value.String())
	}
	if !isInt {
		u64 = uint64(i64)
	}
	return u64, value.String(), nil
}
