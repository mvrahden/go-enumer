package common

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mvrahden/go-enumer/config"
	"github.com/mvrahden/go-enumer/pkg/utils/slices"
)

type EnumTypeConfig struct {
	Node *ast.Comment

	*config.Options
	FromSource string
}

func DefaultConfig(cfg *config.Options) *EnumTypeConfig {
	return &EnumTypeConfig{Options: cfg.Clone()}
}

type EnumType struct {
	Node       *ast.GenDecl
	Config     *EnumTypeConfig
	Spec       *EnumTypeSpec // hint: the specification derived either from const block notation or from a file
	ConstBlock *EnumConstBlock
}

func (e *EnumType) Name() *ast.Ident {
	return e.Node.Specs[0].(*ast.TypeSpec).Name
}

// HasFileSpec indicates whether or not the EnumType is configured
// with a source file. Its usage is legal for AFTER the config has
// been loaded.
func (e *EnumType) HasFileSpec() bool {
	return len(e.Config.FromSource) > 0
}

func (e *EnumType) HasSimpleBlockSpec() bool {
	return !e.HasFileSpec()
}

func (e *EnumType) GetTypeVia(ti *types.Info) types.Type {
	return ti.TypeOf(e.Node.Specs[0].(*ast.TypeSpec).Name)
}

type EnumConstBlock struct {
	Node  *ast.GenDecl
	Specs []*EnumValueSpec // hint: all value specs of a block
}

type EnumValueSpec struct {
	Node  *ast.ValueSpec
	Value uint64
}

func (e *EnumValueSpec) GetTypeVia(ti *types.Info) types.Type {
	return ti.TypeOf(e.Node.Names[0])
}

func (e *EnumValueSpec) GetObjectVia(ti *types.Info) types.Object {
	return ti.ObjectOf(e.Node.Names[0])
}

// DetectMagicComment retrieves the magic comment from the list of comments.
// It assumes that a magic comment exists.
func (e *EnumType) DetectMagicComment() (c *ast.Comment) {
	mcIdx := slices.FindIndex(e.Node.Doc.List, func(c *ast.Comment, idx int) bool {
		return MAGIC_MARKER.MatchString(c.Text)
	})
	return e.Node.Doc.List[mcIdx]
}

// ExtractCommentString is a Noop func, but allows to be intercepted during tests
// to clean comment texts from test artifacts while keeping production code unaffected.
var ExtractCommentString = func(c *ast.Comment) string {
	return c.Text
}

func (e *EnumType) ParseMagicComment(mc *ast.Comment, opts *config.Options) error {
	doc := ExtractCommentString(mc)

	cfg := DefaultConfig(opts)

	if args := strings.Split(doc, " "); len(args) > 1 {
		args = args[1:] /* hint: parse w/o magic marker */
		var f flag.FlagSet
		f.SetOutput(io.Discard) // silence flagset StdErr output

		f.StringVar(&cfg.TransformStrategy, "transform", cfg.TransformStrategy, "")
		f.Var(&cfg.Serializers, "serializers", "")
		f.Var(&cfg.SupportedFeatures, "support", "")
		f.StringVar(&cfg.FromSource, "from", "", "")
		err := f.Parse(args)
		if err != nil {
			if els := strings.SplitAfter(err.Error(), "not defined: -"); len(els) == 2 { // flag provided but not defined: -<unknown opt>
				return fmt.Errorf("unknown option %q", els[1])
			}
			return fmt.Errorf("failed parsing doc-string. err: %s", err)
		}
		if f.NArg() > 0 {
			// report non-flag arguments
			return fmt.Errorf("unknown args %v", f.Args())
		}
	}

	if len(cfg.FromSource) > 0 {
		if !strings.HasSuffix(cfg.FromSource, ".csv") {
			return errors.New("unsupported file extension")
		}
		if strings.Contains(cfg.FromSource, "../") {
			return errors.New("source path cannot contain path traversals")
		}
		if strings.HasPrefix(cfg.FromSource, "./") || strings.HasPrefix(cfg.FromSource, "/") {
			return errors.New("source path cannot start with \"./\" or \"/\"")
		}
		cfg.FromSource = filepath.Clean(cfg.FromSource)
	}

	e.Config = cfg
	e.Config.Node = mc

	return nil
}

func (e *EnumType) GetPkgFS(fset *token.FileSet) (fs.FS, bool) {
	if e.HasSimpleBlockSpec() {
		return nil, false
	}
	dirPath := filepath.Dir(fset.Position(e.Node.Pos()).Filename)
	pkgFs := os.DirFS(dirPath)
	return pkgFs, true
}

func (e *EnumType) ValidateEnumTypeConfig(fset *token.FileSet) error {
	// valdidate simple enum options
	// TODO

	// validate filebased enum options
	pkgFS, ok := e.GetPkgFS(fset)
	if !ok {
		return nil
	}

	_, err := fs.Stat(pkgFS, e.Config.FromSource)
	if errors.Is(err, os.ErrNotExist) {
		return errors.New("no such file")
	} else if err != nil {
		return fmt.Errorf("please verify the file. err: %w", err)
	}
	return nil
}

func (e *EnumType) LoadSpec(fset *token.FileSet) error {
	if e.HasSimpleBlockSpec() {
		return e.LoadSimpleBlockSpec()
	}
	return e.LoadFileSpec(fset)
}

func (e *EnumType) LoadSimpleBlockSpec() error {
	spec := &EnumTypeSpec{Type: SimpleBlockSpec, Values: make([]*EnumTypeSpecValue, len(e.ConstBlock.Specs))}

	slices.Range(e.ConstBlock.Specs, func(v *EnumValueSpec, idx int) {
		enumValue := v.Node.Names[0].Name
		enumValue = strings.TrimPrefix(enumValue, e.Name().Name)
		spec.Values[idx] = &EnumTypeSpecValue{ID: v.Value, EnumValue: enumValue, ConstSpec: v}
	})

	e.Spec = e.detectAlternativeValues(spec)
	return nil
}

func (e *EnumType) LoadFileSpec(fset *token.FileSet) error {
	pkgFS, ok := e.GetPkgFS(fset)
	if !ok {
		return nil
	}
	spec, err := e.loadSpecFromFS(pkgFS)
	if err != nil {
		return err
	}
	e.Spec = e.detectAlternativeValues(spec)
	return nil
}

func (e *EnumType) detectAlternativeValues(spec *EnumTypeSpec) *EnumTypeSpec {
	slices.Range(spec.Values, func(v *EnumTypeSpecValue, idx int) {
		if idx == 0 {
			return
		}
		prev := spec.Values[idx-1]
		v.IsAlternative = prev.ID == v.ID
	})
	return spec
}

const maxSize int64 = 5e6 // 5MB

func (e *EnumType) loadSpecFromFS(pkgFS fs.FS) (*EnumTypeSpec, error) {
	f, err := pkgFS.Open(e.Config.FromSource)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	{ // inspect file
		stat, err := f.Stat()
		if err != nil {
			return nil, err
		}
		if stat.Size() > maxSize {
			return nil, fmt.Errorf("filesize exceeds maximum threshold of 5MB")
		}
	}
	// parse file
	spec := EnumTypeSpec{Type: FilebasedSpec}
	{
		cr := csv.NewReader(f)
		{ // evaluate header
			hdr, err := cr.Read()
			if errors.Is(err, io.EOF) {
				return nil, errors.New("found empty csv source")
			}
			if err != nil {
				return nil, fmt.Errorf("failed reading header. err: %w", err)
			}
			ok := len(hdr) >= 2
			if !ok {
				return nil, errors.New("header must contain at least 2 fields")
			}
			ok = slices.None(hdr, func(v string, _ int) bool {
				return v == ""
			})
			if !ok {
				return nil, errors.New("header cannot contain empty fields")
			}
			ok = slices.None(hdr, func(v string, _ int) bool {
				return IS_NUMERIC_VALUE.MatchString(v)
			})
			if !ok {
				return nil, errors.New("header cannot contain numeric values")
			}
			additionalDataColumns := len(hdr) - 2
			if additionalDataColumns > 0 {
				spec.AdditionalData = &AdditionalData{
					Headers: make([]*AdditionalDataHeader, additionalDataColumns),
				}
				_, err = slices.RangeErr(hdr[2:], func(cell string, idx int) error {
					isTyped := IS_TYPED_HEADER.MatchString(cell)
					if !isTyped { // ok, treat is as string type
						spec.AdditionalData.Headers[idx] = &AdditionalDataHeader{Name: cell, Type: types.String}
						return nil
					}
					// detect type from pattern
					openIdx := strings.Index(cell, "(")
					rawTypeValue := cell[:openIdx]
					typ, ok := getTypeFromString(rawTypeValue)
					if !ok {
						return errors.New("header types can only be native types")
					}
					cellString := cell[openIdx+1 : len(cell)-1]
					spec.AdditionalData.Headers[idx] = &AdditionalDataHeader{Name: cellString, Type: typ}
					return nil
				})
				if err != nil {
					return nil, err
				}
			}
		}
		var fieldValueFormatter = func(t types.BasicKind, val string) string {
			if t == types.String {
				return strconv.Quote(val)
			}
			return val
		}
		{ // evaluate rows
			for rowIdx := 0; true; rowIdx++ {
				row, err := cr.Read()
				if errors.Is(err, io.EOF) {
					if rowIdx == 0 {
						return nil, fmt.Errorf("csv source must contain at least one value row")
					}
					break
				}
				if errors.Is(err, csv.ErrFieldCount) {
					return nil, fmt.Errorf("rows must have same column count as header (see row %d)", rowIdx+2)
				}
				if err != nil {
					return nil, err
				}
				if len(row) < 2 {
					return nil, fmt.Errorf("rows must contain at least 2 columns (see row %d)", rowIdx+2)
				}
				id, err := strconv.ParseUint(row[0], 10, 64)
				if err != nil {
					return nil, fmt.Errorf("failed converting %q to uint64", row[0])
				}
				if rowIdx == 0 && id > 1 {
					return nil, errors.New("enum sequences must start with either 0 or 1")
				}
				val := row[1]
				spec.Values = append(spec.Values, &EnumTypeSpecValue{ID: id, EnumValue: val})

				if spec.AdditionalData == nil {
					continue // no additional data, let's move to next row
				}
				dataCells := row[2:]
				// add a row of additional data
				spec.AdditionalData.Rows = append(spec.AdditionalData.Rows, make([]*AdditionalDataCell, len(dataCells)))
				dataRowIdx := len(spec.AdditionalData.Rows) - 1
				// parse and format additional data
				badColIdx, err := slices.RangeErr(dataCells, func(v string, colIdx int) error {
					typ := spec.AdditionalData.Headers[colIdx].Type
					litVal := fieldValueFormatter(typ, v)
					typedVal, err := typedParserFuncs[typ](v)
					if err != nil {
						// hint: float special types are set to "0" here
						// TODO: a proper implementation would map to "math.Inf" and "math.NaN" funcs
						// and add "math" package import to file
						if errors.Is(err, ErrIsNaN) || errors.Is(err, ErrIsPosInf) || errors.Is(err, ErrIsNegInf) {
							litVal = "0"
							goto CONTINUE_ASSIGNMENT
						}
						return err
					}
				CONTINUE_ASSIGNMENT:
					spec.AdditionalData.Rows[dataRowIdx][colIdx] = &AdditionalDataCell{LiteralValue: litVal, TypedValue: typedVal}
					return nil
				})
				if err != nil {
					return nil, fmt.Errorf("failed parsing additional data in row %d column %d. err: %w", rowIdx+2, badColIdx+2, err)
				}
			}
		}
	}
	return &spec, nil
}

func (e *EnumType) ValidateConstBlock(fset *token.FileSet, typesInfo *types.Info) error {
	if e.ConstBlock == nil {
		if e.HasSimpleBlockSpec() {
			return errors.New("enum types require a const block or a file source")
		}
		return nil
	}

	// assert const block is in same file
	if fset.File(e.ConstBlock.Node.Pos()) != fset.File(e.Node.Pos()) {
		return errors.New("enum const block must be in same file as their type definition")
	}
	// assert const block is after relevant type
	if e.ConstBlock.Node.Pos() < e.Node.Pos() {
		return errors.New("enum const block must be defined after its type definition")
	}
	// assert const block has no rowed declarations
	ok := slices.None(e.ConstBlock.Specs, func(v *EnumValueSpec, idx int) bool {
		return len(v.Node.Names) > 1 || len(v.Node.Values) > 1
	})
	if !ok {
		return errors.New("enum const block must not contain rowed declarations")
	}
	// assert only enum values of relevant enum type within block
	ok = slices.None(e.ConstBlock.Specs, func(curr *EnumValueSpec, idx int) bool {
		if idx == 0 { // assert that first type declaration is an enum
			isSameAsBlockType := types.IdenticalIgnoreTags(
				curr.GetTypeVia(typesInfo),
				e.GetTypeVia(typesInfo),
			)
			return !isSameAsBlockType
		}
		relevantEnumType := e.ConstBlock.Specs[0].GetTypeVia(typesInfo)
		isSameAsBlockType := types.IdenticalIgnoreTags(
			typesInfo.TypeOf(curr.Node.Names[0]),
			relevantEnumType,
		)
		return !isSameAsBlockType
	})
	if !ok {
		return errors.New("enum const block must not contain unrelated type declarations")
	}

	// assert no skipped rows in blocks
	ok = slices.None(e.ConstBlock.Specs, func(vs *EnumValueSpec, idx int) bool {
		// Special Case: "skipped rows" provoke an increment of more than one.
		return vs.Node.Names[0].Name == "_"
	})
	if !ok {
		return errors.New("enum const block must not contain skipped rows")
	}

	// assert numerical correctness
	ok = slices.All(e.ConstBlock.Specs, func(vs *EnumValueSpec, idx int) bool {
		val := vs.GetObjectVia(typesInfo).(*types.Const).Val()
		{
			val, ok := constant.Int64Val(val)
			if !ok {
				return false
			}
			// hint: reflow otherwise overflown values.
			// that's ok as we're dealing with uint types exclusively
			vs.Value = uint64(val)
		}
		return true
	})
	if !ok {
		return errors.New("invalid numerical format")
	}
	return nil
}

func (e *EnumType) ValidateSpec(fset *token.FileSet, typesInfo *types.Info) error {
	if len(e.Spec.Values) == 0 {
		return errors.New("enum spec must contain at least one value")
	}

	// assert spec sequence start
	if e.Spec.Values[0].ID > 1 {
		// hint: file based enums can start with arbitrary numbers in const blocks
		// as they do not represent the SPEC in this case but merely refer to individual values.
		return errors.New("enum spec sequences must start with either 0 or 1")
	}

	// assert order of values
	ok := slices.None(e.Spec.Values, func(v *EnumTypeSpecValue, idx int) bool {
		if idx == 0 {
			return false
		}
		prev := e.Spec.Values[idx-1].ID
		return prev > v.ID
	})
	if !ok {
		return errors.New("enum spec sequences must be ordered")
	}

	// assert increments of values
	ok = slices.None(e.Spec.Values, func(v *EnumTypeSpecValue, idx int) bool {
		if idx == 0 {
			return false
		}
		prev := e.Spec.Values[idx-1].ID
		return prev+1 < v.ID
	})
	if !ok {
		return errors.New("enum spec sequences must increment at most by one")
	}
	return nil
}

func (e *EnumType) CrossValidateConstBlockWithSpec(fset *token.FileSet, typesInfo *types.Info) error {
	if e.ConstBlock == nil {
		return nil
	}
	if e.HasSimpleBlockSpec() {
		// block specs are derived from blocks and have gone through
		// extensive assertions
		return nil
	}
	// Enums based on file specs can have additional const blocks
	// to reference individual often used values for convenience.
	// Here we need to ensure that these const block values are:
	// - complying to good usage practices
	// - in sync with their spec

	// assert one value reference per enum type only
	badIdx := slices.FindIndex(e.ConstBlock.Specs, func(s1 *EnumValueSpec, idx int) bool {
		if idx == len(e.ConstBlock.Specs)-1 {
			return false
		}
		return slices.Any(e.ConstBlock.Specs[idx+1:], func(s2 *EnumValueSpec, idx int) bool {
			return s1.Value == s2.Value
		})
	})
	if badIdx > -1 {
		constValue := e.ConstBlock.Node.Specs[badIdx].(*ast.ValueSpec)
		return fmt.Errorf("%q is a redundant constant", constValue.Names[0].Name)
	}

	// assert const block values do not exceed spec range
	specMin := e.Spec.Values[0].ID
	specMax := e.Spec.Values[len(e.Spec.Values)-1].ID
	badIdx = slices.FindIndex(e.ConstBlock.Specs, func(vs *EnumValueSpec, idx int) bool {
		return vs.Value < specMin || vs.Value > specMax
	})
	if badIdx > -1 {
		constValue := e.ConstBlock.Node.Specs[badIdx].(*ast.ValueSpec)
		return fmt.Errorf("%q exceeds spec range [%d,%d]", constValue.Names[0].Name, specMin, specMax)
	}

	// assert const block assertions
	badIdx, err := slices.RangeErr(e.ConstBlock.Specs, func(vs *EnumValueSpec, _ int) error {
		if vs.Node.Comment != nil {
			assertToken := "// assert \""
			lineComment := vs.Node.Comment.List[0].Text
			if strings.HasPrefix(lineComment, assertToken) {
				endIdx := strings.Index(lineComment[len(assertToken):], "\"")
				if endIdx == -1 {
					return errors.New("missing terminating quote in assertion")
				}
				assertVal := lineComment[len(assertToken) : len(assertToken)+endIdx]

				specIdx := slices.FindIndex(e.Spec.Values, func(v *EnumTypeSpecValue, idx int) bool {
					return v.ID == vs.Value
				})
				if e.Spec.Values[specIdx].EnumValue != assertVal {
					return errors.New("assertion failed")
				}
			}
		}
		return nil
	})
	if badIdx > -1 {
		constValue := e.ConstBlock.Node.Specs[badIdx].(*ast.ValueSpec)
		return fmt.Errorf("%q fails on assertion (reason: %s)", constValue.Names[0].Name, err)
	}
	return nil
}
