package gen

import (
	"bytes"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/mvrahden/go-enumer/config"
	"github.com/mvrahden/go-enumer/pkg/utils/slices"
	"golang.org/x/tools/go/packages"
)

const MAGIC_MARKER = "//go:enum"

var (
	matchGeneratedFileRegex = regexp.MustCompile(`^// Code generated .* DO NOT EDIT.$`)
	matchNumericValueRegex  = regexp.MustCompile(`^\-?\d+`)
	matchTypedHeaderRegex   = regexp.MustCompile(`^.+(\(.+\))$`)
)

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
	if err := i.inspectTypeSpecs(pkg, out); err != nil {
		return nil, err
	}
	sources, err := i.inspectDocstrings(pkg, out)
	if err != nil {
		return nil, err
	}
	if err := i.inspectSources(out, sources); err != nil {
		return nil, err
	}
	i.inspectImports(out)

	i.sortTypeSpecs(out)
	if err := i.validateFile(out); err != nil {
		return nil, err
	}
	return out, nil
}

func (i inspector) sortTypeSpecs(f *File) {
	// sort all enums and all enum values
	f.TypeSpecs = slices.SortStable(f.TypeSpecs, func(s []*TypeSpec, i, j int) bool {
		return strings.Compare(s[i].Name, s[j].Name) < 0
	})
	slices.Range(f.TypeSpecs, func(v *TypeSpec, idx int) {
		v.ValueSpecs = slices.SortStable(v.ValueSpecs, func(s []*ValueSpec, i, j int) bool {
			return s[i].Value < s[j].Value
		})
	})
}

func (i inspector) inspectDocstrings(pkg *packages.Package, f *File) (sources []string, err error) {
	sources = make([]string, len(f.TypeSpecs))
	for idx, ts := range f.TypeSpecs {
		args := strings.Split(ts.Meta.Docstring, " ")
		ts.Meta.Config = i.cfg.Clone()

		if len(args) <= 1 {
			continue
		}

		var fs flag.FlagSet
		fs.SetOutput(io.Discard) // silence flagset StdErr output

		fs.StringVar(&ts.Meta.Config.TransformStrategy, "transform", ts.Meta.Config.TransformStrategy, "")
		fs.Var(&ts.Meta.Config.Serializers, "serializers", "")
		fs.Var(&ts.Meta.Config.SupportedFeatures, "support", "")
		fs.StringVar(&sources[idx], "from", "", "")
		err := fs.Parse(args[1:]) // hint: parse w/o magic marker
		if err != nil {
			return nil, fmt.Errorf("Failed parsing doc-string for %q. err: %w", ts.Name, err)
		}

		if len(sources[idx]) == 0 {
			continue
		}
		if strings.Contains(sources[idx], "../") {
			return nil, fmt.Errorf("Invalid source file path in doc-string for %q. err: forbidden path traversal detected", ts.Name)
		}
		if strings.HasPrefix(sources[idx], "./") {
			return nil, fmt.Errorf("Invalid source file path in doc-string for %q. err: cannot start with \"./\"", ts.Name)
		}
	}
	return sources, nil
}

func (i inspector) inspectSources(f *File, sources []string) error {
	for idx, ts := range f.TypeSpecs {
		fromSource := sources[idx]
		if len(fromSource) == 0 {
			continue
		}

		err := i.readFromCSV(ts, fromSource)
		if err != nil {
			return fmt.Errorf("Failed reading from CSV for %q. err: %w", ts.Name, err)
		}
	}
	return nil
}

func (i inspector) readFromCSV(ts *TypeSpec, p string) error {
	buf, err := os.ReadFile(filepath.Join(filepath.Dir(ts.Meta.Filepath), p))
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("no such file %q", p)
	} else if errors.Is(err, os.ErrPermission) {
		return fmt.Errorf("no permission for file %q", p)
	} else if err != nil {
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
	ts.Meta.IsFromCsvSource = true // hint: mark type as derived from CSV

	if len(records) == 0 {
		return fmt.Errorf("csv source must contain a header row")
	}

	// trim whitespace on all cells
	for _, row := range records {
		for idx := range row {
			row[idx] = strings.TrimSpace(row[idx])
		}
	}

	// evaluate header row
	for _, row := range records[0:1] {
		ok := matchNumericValueRegex.MatchString(row[0])
		if ok {
			return fmt.Errorf("first row must be a header row but found numeric value in first cell")
		}

		if len(row) > 2 { // add additional header names
			ts.Meta.HasAdditionalData = true
		}
	}

	// evaluate additional data of header row
	for _, row := range records[0:1] {
		if !ts.Meta.HasAdditionalData {
			break
		}

		ts.AdditionalData = &AdditionalData{
			Headers: make([]AdditionalDataHeader, len(row[2:])),
		}

		for idx, cell := range row[2:] {
			if isUntyped := !matchTypedHeaderRegex.MatchString(cell); isUntyped {
				ts.AdditionalData.Headers[idx] = AdditionalDataHeader{cell, GoTypeUnknown}
				continue
			}
			// determine header type
			openIdx := strings.Index(cell, "(")
			typeValue := cell[:openIdx]
			cellString := cell[openIdx+1 : len(cell)-1]
			determinedType := getTypeFromString(typeValue)
			ts.AdditionalData.Headers[idx] = AdditionalDataHeader{cellString, determinedType}
		}
	}

	// drop header row after evaluation
	records = records[1:]

	if len(records) == 0 {
		return fmt.Errorf("csv source must contain at least one value row")
	}

	// evaluate data rows (base data)
	csvValuespecs := make([]*ValueSpec, len(records))
	for idx, row := range records {
		rowId, err := strconv.ParseUint(row[0], 10, 64)
		if err != nil {
			return fmt.Errorf("failed converting %q to uint64 at line %d", row[0], idx+2)
		}
		if idx == 0 && rowId > 1 {
			return fmt.Errorf("found invalid start of enum sequence at line %d.", idx+2)
		}
		if idx > 0 && rowId < csvValuespecs[idx-1].Value {
			return fmt.Errorf("found decreasing value at line %d.", idx+2)
		}

		csvValuespecs[idx] = &ValueSpec{
			Value:              rowId,
			String:             row[1],
			ConstName:          "", // hint: no identifier here
			IsAlternativeValue: idx != 0 && csvValuespecs[idx-1].Value == rowId,
		}
	}

	{
		// cross-validate value specs
		// enum constants must be within extent of csv values
		min, max := csvValuespecs[0].Value, csvValuespecs[len(csvValuespecs)-1].Value
		if idx := slices.FindIndex(ts.ValueSpecs, func(v *ValueSpec, idx int) bool {
			return v.Value < min || v.Value > max
		}); idx > -1 {
			return fmt.Errorf("enum constant %q is out of csv source range [%d,%d]", ts.ValueSpecs[idx].ConstName, min, max)
		}
	}

	// swap enum constants with csv value specs
	ts.ValueSpecs = csvValuespecs

	// evaluate data rows (additional data)
	if !ts.Meta.HasAdditionalData {
		return nil
	}
	// reduce dataset to additional data only
	records = slices.Map(records, func(row []string, idx int) []string {
		return row[2:] // drop first two columns
	})

	// assert that alternative data is either zero all identical to master value
	if idx := slices.FindIndex(records, func(row []string, rowIdx int) bool {
		if !csvValuespecs[rowIdx].IsAlternativeValue {
			return false // hint: we only care for alternative values
		}
		if ok := slices.All(row, func(v string, colIdx int) bool {
			return len(v) == 0
		}); ok {
			return false // hint: all fields are empty, that's ok
		}
		if rowIdx == 0 {
			return false
		}
		return slices.Any(row, func(v string, colIdx int) bool {
			// hint: all fields must be identical to master value
			return records[rowIdx-1][colIdx] != v
		})
	}); idx != -1 {
		return fmt.Errorf("found invalid additional data of an alternative enum value %q at line %d.", csvValuespecs[idx].String, idx+2)
	}

	// filter relevant records
	records = slices.Filter(records, func(v []string, idx int) bool {
		// hint: alternative values can not have deviating additional data
		// therefore we must filter them
		return !csvValuespecs[idx].IsAlternativeValue
	})

	// prepare parsers and type formatters
	columnParseFuncs := make([]func(raw string) (any, error), len(ts.AdditionalData.Headers))
	columnTypeFormatter := make([]func(raw any) string, len(ts.AdditionalData.Headers))
	for idx, col := range ts.AdditionalData.Headers {
		columnParseFuncs[idx] = getParserFuncFor(col.Type)
		columnTypeFormatter[idx] = col.Type.ToSource
	}

	ts.AdditionalData.Rows = make([][]AdditionalDataCell, len(records))
	for rowIdx, row := range records {
		ts.AdditionalData.Rows[rowIdx] = make([]AdditionalDataCell, len(ts.AdditionalData.Headers))
		for colIdx, raw := range row {
			literalValue := raw
			parseByColumnType := columnParseFuncs[colIdx]
			formatByColumnType := columnTypeFormatter[colIdx]
			cellValue, err := parseByColumnType(raw)
			// hint: float special types are set to "0" here
			// a proper implementation would map to "math.Inf" and "math.NaN" funcs
			// and add "math" package import to file
			if errors.Is(err, ErrIsNaN) {
				literalValue = "0"
				err = nil
			} else if errors.Is(err, ErrIsPosInf) {
				literalValue = "0"
				err = nil
			} else if errors.Is(err, ErrIsNegInf) {
				literalValue = "0"
				err = nil
			}
			if err != nil {
				return err
			}
			ts.AdditionalData.Rows[rowIdx][colIdx] = AdditionalDataCell{
				LiteralValue: formatByColumnType(literalValue),
				RawValue:     cellValue,
			}
		}
	}

	return nil
}

func (i inspector) inspectImports(f *File) {
	f.Imports = append(f.Imports, &Import{Path: "errors"})
	f.Imports = append(f.Imports, &Import{Path: "fmt"})

	// we add all imports (also duplicates)
	for _, ts := range f.TypeSpecs {
		for _, v := range ts.Meta.Config.Serializers {
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
	// sort alphabetically
	f.Imports = slices.SortStable(f.Imports, func(s []*Import, i, j int) bool {
		return strings.Compare(s[i].Path, s[j].Path) < 0 && strings.Compare(s[i].Name, s[j].Name) < 0
	})
	// filter unique import values
	f.Imports = slices.Filter(f.Imports, func(v *Import, idx int) bool {
		if idx == 0 {
			return true
		}
		return v.Path != f.Imports[idx-1].Path || v.Name != f.Imports[idx-1].Name
	})
}

func (i inspector) validateFile(f *File) error {
	for _, ts := range f.TypeSpecs {
		// validate start value of enum sequence
		if len(ts.ValueSpecs) > 0 &&
			!(ts.ValueSpecs[0].Value == 0 || ts.ValueSpecs[0].Value == 1) {
			return fmt.Errorf("Enum %q must start with either 0 or 1.", ts.Name)
		}
		// ensure we have a linearly incrementing sequence of values.
		// however, an enum can assign a numeric value multiple times.
		// therefore we must only dismiss distances > 1.
		for idx := 1; idx < len(ts.ValueSpecs); idx++ {
			delta := ts.ValueSpecs[idx].Value - ts.ValueSpecs[idx-1].Value
			if delta > 1 {
				return fmt.Errorf("Enum %q must be a continuous sequence with linear increments of 1.", ts.Name)
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

func (i inspector) inspectTypeSpecs(pkg *packages.Package, out *File) error {
	specs, err := i.determinePackageScopedEnumTypeSpecs(pkg, out)
	if err != nil {
		return err
	}
	for idx, s := range specs {
		typ, err := i.determineTypeOfExpr(s.TypeSpec.Name)
		if err != nil {
			return fmt.Errorf("Enum type of %q %w", s.TypeSpec.Name, err)
		}
		out.TypeSpecs = append(out.TypeSpecs, &TypeSpec{
			Name: s.TypeSpec.Name.Name,
			Type: typ,
			Meta: &MetaTypeSpec{
				Docstring: s.EnumMarker,
				Filepath:  s.File.Name(),
			},
		})
		vspecs, err := i.evaluateValueSpecs(s, pkg)
		if err != nil {
			return err
		}
		out.TypeSpecs[idx].ValueSpecs = vspecs
	}
	return nil
}

// Naive check whether file contains a "Code generated" comment.
// This function does not verify the position of the comment.
func (i inspector) isGeneratedFile(f *ast.File) bool {
	if len(f.Comments) > 0 &&
		len(f.Comments[0].List) > 0 {
		firstComment := f.Comments[0].List[0]
		if matchGeneratedFileRegex.MatchString(firstComment.Text) {
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
	Values     []*ast.Ident
}

func (i *inspector) determinePackageScopedEnumTypeSpecs(pkg *packages.Package, out *File) ([]typeSpec, error) {
	var typeSpecs []typeSpec
	for _, f := range pkg.Syntax {
		// determine type specs
		for _, v := range f.Decls {
			decl, ok := v.(*ast.GenDecl)
			if !ok || decl.Tok != token.TYPE {
				// we only care about type declarations.
				continue
			}
			if len(decl.Specs) != 1 {
				continue
			}

			// Detect magic comment in doc-string
			if decl.Doc == nil || len(decl.Doc.List) == 0 {
				continue // missing doc-string
			}
			var magicComment string
			{
				lastDoc := decl.Doc.List[len(decl.Doc.List)-1]
				docString := strings.TrimSpace(lastDoc.Text)
				hasMagicComment := strings.Compare(docString, MAGIC_MARKER) == 0 ||
					strings.HasPrefix(docString, MAGIC_MARKER+" ") // hint: w/o OR w/ subsequent config string

				if !hasMagicComment {
					continue // no magic comment
				}
				magicComment = docString
			}

			ts, ok := decl.Specs[0].(*ast.TypeSpec)
			if !ok {
				continue // not a type spec
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
			typeSpecs = append(typeSpecs, typeSpec{Decl: decl, TypeSpec: ts, Type: btyp, File: f, EnumMarker: magicComment})
		}

		// determine values for type specs
		for _, v := range f.Decls {
			decl, ok := v.(*ast.GenDecl)
			if !ok || decl.Tok != token.CONST {
				// we only care about const declarations.
				continue
			}

			var prevType *ast.Ident // needed for blocks with implicit types
			for _, spec := range decl.Specs {
				vspec, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				if vspec.Type == nil && prevType == nil {
					// no type information available
					continue
				}

				if vspec.Type == nil && prevType != nil {
					goto HAS_TYPE
				}
				if ident, ok := vspec.Type.(*ast.Ident); !ok {
					// not the type we're searching for
					continue
				} else if ident == nil && prevType == nil {
					continue
				} else if ident != nil {
					prevType = ident
				}
			HAS_TYPE:

				if vspec.Type == nil {
					// patch missing type information for those specs
					// with implicit type
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
							typeSpecs[idx].Values = append(typeSpecs[idx].Values, v)
						}
					}
				}
			}
		}
	}
	return typeSpecs, nil
}

func (i inspector) evaluateValueSpecs(t typeSpec, pkg *packages.Package) ([]*ValueSpec, error) {
	vspecs := make([]*ValueSpec, len(t.Values))
	for idx, value := range t.Values {
		val, err := i.determineValueOfExpr(value, pkg)
		if err != nil {
			return nil, err
		}
		name := value.Name
		// auto-strip prefix
		if strings.HasPrefix(value.Name, t.TypeSpec.Name.Name) {
			if len(t.TypeSpec.Name.Name) == len(value.Name) {
				return nil, fmt.Errorf("failed to auto-strip prefix (enum value equals type name). make sure to give a meaningful names to your enum values.")
			}
			name = value.Name[len(t.TypeSpec.Name.Name):]
		}
		vspecs[idx] = &ValueSpec{
			ConstName:          value.Name,
			String:             name,
			Value:              val,
			IsAlternativeValue: idx != 0 && val == vspecs[idx-1].Value,
		}
	}
	return vspecs, nil
}

func (i inspector) determineTypeOfExpr(e ast.Expr) (GoType, error) {
	switch t := e.(type) {
	case *ast.Ident:
		typ, ok := validEnumTypesMap[t.Name]
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
	return GoTypeUnknown, fmt.Errorf("must be of an unsigned integer type, found %q.", e)
}

func (i inspector) determineValueOfExpr(e ast.Expr, pkg *packages.Package) (uint64, error) {
	c, ok := e.(*ast.Ident)
	if !ok {
		return 0, fmt.Errorf("internal error: a value slipped our type evaluation (type: %+v)", e)
	}
	obj := pkg.TypesInfo.ObjectOf(c)
	if obj == nil {
		return 0, fmt.Errorf("no type object for constant %q", c)
	}
	objT := obj.Type()
	if objT == nil {
		return 0, fmt.Errorf("definition type for constant %q is <nil>", c)
	}
	ul := objT.Underlying()
	if ul == nil {
		return 0, fmt.Errorf("underlying was expected to be a basic type, but was <nil>")
	}
	bul, ok := ul.(*types.Basic)
	if !ok {
		return 0, fmt.Errorf("underlying type for constant %q is not a basic type", c)
	}
	info := bul.Info()
	if info&types.IsInteger == 0 {
		return 0, fmt.Errorf("%q is not a constant of type integer", c.Name)
	}
	cobj, ok := obj.(*types.Const)
	if !ok {
		return 0, fmt.Errorf("internal error: a value slipped our type evaluation (type: %+v is not a const)", e)
	}
	value := cobj.Val()
	if value.Kind() != constant.Int {
		return 0, fmt.Errorf("constant is not an integer %q", c)
	}
	u64, isUint := constant.Uint64Val(value)
	if !isUint {
		return 0, fmt.Errorf("internal error: value of %q is not an unsigned integer: %s", c, value)
	}
	return u64, nil
}
