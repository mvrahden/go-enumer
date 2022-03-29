package gen

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
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
	if err := i.inspectDocstrings(pkg, out); err != nil {
		return nil, err
	}

	if err := i.validateFile(out); err != nil {
		return nil, err
	}
	i.inspectImports(out)
	return out, nil
}

func (i inspector) inspectDocstrings(pkg *packages.Package, f *File) error {
	for _, ts := range f.TypeSpecs {
		args := strings.Split(ts.Docstring, " ")
		newCfg := i.cfg.Clone()
		if len(args) <= 1 {
			ts.Config = *newCfg
			continue
		}
		var fs flag.FlagSet
		var fromSource string
		fs.StringVar(&newCfg.TransformStrategy, "transform", newCfg.TransformStrategy, "")
		fs.Var(&newCfg.Serializers, "serializers", "")
		fs.Var(&newCfg.SupportedFeatures, "support", "")
		fs.StringVar(&fromSource, "from", "", "")
		err := fs.Parse(args[1:])
		if err != nil {
			return err
		}

		if len(fromSource) == 0 {
			continue
		}
		err = i.readFromCSV(ts, fromSource)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i inspector) readFromCSV(ts *TypeSpec, p string) error {
	dir := filepath.Dir(ts.Filepath)
	p = filepath.Join(dir, p)
	buf, err := os.ReadFile(p)
	if err != nil {
		return err
	}
	if len(buf) == 0 {
		return fmt.Errorf("found empty csv source")
	}
	cr := csv.NewReader(bytes.NewBuffer(buf))
	records, err := cr.ReadAll()
	if err != nil {
		return err
	}
	ts.IsFromCsvSource = true    // hint: mark type as derived from CSV
	ts.HasCanonicalValues = true // hint: mark type for canocical value support
	ts.ValueSpecs = make([]*ValueSpec, len(records))
	for idx, row := range records {
		u64, err := strconv.ParseUint(row[0], 10, 64)
		if err != nil {
			return fmt.Errorf("failed converting %q to uint64", row[0])
		}
		ts.ValueSpecs[idx] = &ValueSpec{
			Index:          idx,
			Value:          u64,
			IdentifierName: "", // hint: no identifier here
			EnumValue:      row[1],
			ValueString:    row[0],
		}
		if len(row) == 3 {
			ts.ValueSpecs[idx].CanonicalValue = row[2]
		}
	}
	return nil
}

func (i inspector) inspectImports(f *File) {
	f.Imports = append(f.Imports, &Import{Path: "fmt"})
	for _, v := range i.cfg.Serializers {
		switch v {
		case "gql":
			f.Imports = append(f.Imports, &Import{Path: "io"})
			f.Imports = append(f.Imports, &Import{Path: "strconv"})
		case "json":
			f.Imports = append(f.Imports, &Import{Path: "encoding/json"})
		case "sql":
			f.Imports = append(f.Imports, &Import{Path: "database/sql/driver"})
		case "yaml.v3":
			f.Imports = append(f.Imports, &Import{Path: "gopkg.in/yaml.v3"})
		}
	}
}

func (i inspector) validateFile(f *File) error {
	for _, v := range f.TypeSpecs {
		sort.SliceStable(v.ValueSpecs, func(i, j int) bool {
			return v.ValueSpecs[i].Value < v.ValueSpecs[j].Value
		})
		// validate start value of enum sequence
		if len(v.ValueSpecs) > 0 &&
			!(v.ValueSpecs[0].Value == 0 || v.ValueSpecs[0].Value == 1) {
			return fmt.Errorf("Invalid enum set: Enums need to start with either 0 or 1.")
		}
		// ensure we have a linearly incrementing sequence of values.
		// however, an enum can assign a numeric value multiple times.
		// therefore we must only dismiss distances > 1.
		for idx := 1; idx < len(v.ValueSpecs); idx++ {
			delta := v.ValueSpecs[idx].Value - v.ValueSpecs[idx-1].Value
			if delta > 1 {
				return fmt.Errorf("Invalid enum set: Enums must be a continuous sequence with linear increments of 1.")
			}
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
	specs, err := i.determinePackageScopedEnumTypeSpecs(pkg, out)
	if err != nil {
		return err
	}
	for idx, s := range specs {
		typ, err := i.determineTypeOfExpr(s.TypeSpec.Name)
		if err != nil {
			return err
		}
		out.TypeSpecs = append(out.TypeSpecs, &TypeSpec{
			Index:     idx,
			Name:      s.TypeSpec.Name.Name,
			Type:      typ,
			Docstring: s.EnumMarker,
			Filepath:  s.File.Name(),
		})
		for _, v := range s.values {
			vspec, err := i.evaluateValueSpec(s, v, pkg)
			if err != nil {
				return err
			}
			out.TypeSpecs[idx].ValueSpecs = append(out.TypeSpecs[idx].ValueSpecs, vspec)
		}
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

type typeSpec struct {
	Decl       *ast.GenDecl
	TypeSpec   *ast.TypeSpec
	Type       *types.Basic
	File       *token.File
	EnumMarker string
	values     []valueSpec
}

type valueSpec struct {
	Value *ast.Ident
}

func (i *inspector) determinePackageScopedEnumTypeSpecs(pkg *packages.Package, out *File) ([]typeSpec, error) {
	var typeSpecs []typeSpec
	for _, f := range pkg.Syntax {
		// determine typeSpecs
		for _, v := range f.Decls {
			decl, ok := v.(*ast.GenDecl)
			if !ok || decl.Tok != token.TYPE {
				// We only care about type declarations.
				continue
			}
			if len(decl.Specs) != 1 {
				continue
			}

			// Detect enum marker in doc-string
			if decl.Doc == nil {
				continue
			}
			var isEnum bool
			var enumMarker string
			for _, cv := range decl.Doc.List {
				if strings.HasPrefix(cv.Text, "//go:enumer") {
					isEnum = true
					enumMarker = cv.Text
					break
				}
			}
			if !isEnum {
				continue
			}
			ts, ok := decl.Specs[0].(*ast.TypeSpec)
			if !ok {
				continue
			}

			typ := pkg.TypesInfo.TypeOf(ts.Type)
			if typ == nil {
				return nil, fmt.Errorf("definition type for constant %q is <nil>", typ)
			}
			ul := typ.Underlying()
			if ul == nil {
				return nil, fmt.Errorf("underlying was expected to be a basic type, but was <nil>")
			}
			btyp, ok := ul.(*types.Basic)
			if !ok {
				return nil, fmt.Errorf("underlying type for constant %q is not a basic type", ul)
			}
			f := pkg.Fset.File(decl.TokPos)
			typeSpecs = append(typeSpecs, typeSpec{Decl: decl, TypeSpec: ts, Type: btyp, File: f, EnumMarker: enumMarker})
		}

		// determine values for typeSpecs
		for _, v := range f.Decls {
			decl, ok := v.(*ast.GenDecl)
			if !ok || decl.Tok != token.CONST {
				// We only care about const declarations.
				continue
			}

			var prevType *ast.Ident // for blocks with implicit types
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
						// not the type we're searching for
						continue
					}
					prevType = ident
				}
				if prevType == nil {
					continue
				}
				if vspec.Type == nil {
					// for those with implicit type, assign the previous type
					vspec.Type = prevType
				}
				// determine the right typespec
				for idx, vts := range typeSpecs {
					if vts.TypeSpec.Name.Name == prevType.Name {
						// add all names
						for _, v := range vspec.Names {
							if v.Name == "_" {
								// blank identifier, not what we're interested in
								continue
							}
							typeSpecs[idx].values = append(vts.values, valueSpec{v})
						}
					}
				}
			}
		}
	}
	return typeSpecs, nil
}

func (i inspector) evaluateValueSpec(t typeSpec, s valueSpec, pkg *packages.Package) (*ValueSpec, error) {
	val, valStr, err := i.determineValueOfExpr(s.Value, pkg)
	if err != nil {
		return nil, err
	}
	name := s.Value.Name
	// auto-strip prefix
	if strings.HasPrefix(s.Value.Name, t.TypeSpec.Name.Name) {
		if len(t.TypeSpec.Name.Name) == len(s.Value.Name) {
			return nil, fmt.Errorf("failed to auto-strip prefix (enum value equals type name). make sure to give a meaningful names to your enum values.")
		}
		name = s.Value.Name[len(t.TypeSpec.Name.Name):]
	}
	return &ValueSpec{
		IdentifierName: s.Value.Name,
		EnumValue:      name,
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
	obj := pkg.TypesInfo.ObjectOf(c)
	if obj == nil {
		return 0, "", fmt.Errorf("no type object for constant %q", c)
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
		return 0, "", fmt.Errorf("%q is not a constant of type integer", c.Name)
	}
	cobj, ok := obj.(*types.Const)
	if !ok {
		return 0, "", fmt.Errorf("internal error: a value slipped our type evaluation (type: %+v is not a const)", e)
	}
	value := cobj.Val()
	if value.Kind() != constant.Int {
		return 0, "", fmt.Errorf("constant is not an integer %q", c)
	}
	i64, isInt := constant.Int64Val(value)
	u64, isUint := constant.Uint64Val(value)
	if !isInt && !isUint {
		return 0, "", fmt.Errorf("internal error: value of %q is not an integer: %s", c, value.String())
	}
	if isInt && i64 < 0 {
		return 0, "", fmt.Errorf("Invalid enum set: values cannot be in a negative range.")
	}
	if !isInt {
		u64 = uint64(i64)
	}
	return u64, value.String(), nil
}
