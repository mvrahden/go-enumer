package enumer

import (
	"errors"
	"go/ast"
	"go/token"
	"go/types"
	"math"
	"strconv"

	"github.com/mvrahden/go-enumer/pkg/utils/slices"
)

// DetectGeneratedFile determines the generated enumer file.
// If no such file exists, it will return `nil`.
func DetectGeneratedFile(files []*ast.File) (genFile *ast.File) {
	genFileIdx := slices.FindIndex(files, func(f *ast.File, _ int) bool {
		if len(f.Comments) == 0 {
			return false
		}
		t := f.Comments[0].List[0].Text
		return GEN_ENUMER_FILE.MatchString(t)
	})
	if genFileIdx == -1 {
		return
	}
	return files[genFileIdx]
}

// DetermineEnumType evaluates given node for enum types.
// If the given node is not fulfilling the requirements for a possible enum type declaration it
// returns zero values.
// If the given node violates the requirements for a eum declaration it will return an error and the token position.
func DetermineEnumType(node ast.Node, typesInfo *types.Info, genFile *ast.File) (*EnumType, token.Pos, error) {
	decl, ok := node.(*ast.GenDecl)
	if !ok || decl.Tok != token.TYPE {
		// we only care about type declarations
		return nil, -1, nil
	}
	{ // ensure we are not inspecting anything from the generated file
		if genFile != nil && decl.Pos() >= genFile.Pos() && decl.Pos() <= genFile.End() {
			return nil, -1, nil
		}
	}

	{ // assert enum type
		if decl.Doc == nil || len(decl.Doc.List) == 0 {
			return nil, -1, nil
		}

		// assert underlying enum type
		if len(decl.Specs) != 1 {
			// hint: this should not happen
			return nil, -1, nil // not a type spec
		}
		ts, ok := decl.Specs[0].(*ast.TypeSpec)
		if !ok {
			return nil, -1, nil // not a type spec
		}
		typ, ok := typesInfo.TypeOf(ts.Type).(*types.Basic)
		if !ok {
			typ, ok = typesInfo.TypeOf(ts.Type).Underlying().(*types.Basic)
		}
		if !ok || typ.Kind() < types.Uint || typ.Kind() > types.Uint64 {
			// TODO: evaluate if this error return is correct at this point (we still don't know if it is an enum spec)
			return nil, node.Pos(), errors.New("enum types must be of any unsigned integer type") // not a type spec
		}

		// find magic comment
		magic := slices.Filter(decl.Doc.List, func(v *ast.Comment, idx int) bool {
			return MAGIC_MARKER.MatchString(v.Text)
		})
		if len(magic) == 0 {
			return nil, -1, nil
		}
		if len(magic) > 1 {
			return nil, node.Pos(), errors.New("at most one magic comment permitted per enum type")
		}
		// assert magic comment position
		if decl.Doc.List[len(decl.Doc.List)-1] != magic[0] {
			return nil, magic[0].Pos(), errors.New("magic comment must be last row of doc string for enum type")
		}
	}

	return &EnumType{Node: decl}, -1, nil
}

// AssignEnumConstBlockToType evaluates the current node for a possible const block spec/enum values.
// If the given node does not fulfill the requirements for a possible const block spec, it returns zero values.
// If the given node violates the requirements for possible const block spec, it returns an error and the token position.
// It otherwise determines all const block values and assigns them to relevant enumtype from the given slice of enum types.
func AssignEnumConstBlockToType(node ast.Node, typesInfo *types.Info, genFile *ast.File, enumTypes []*EnumType) (token.Pos, error) {
	decl, ok := node.(*ast.GenDecl)
	if !ok || decl.Tok != token.CONST {
		// we only care about const declarations
		return -1, nil
	}
	{ // ensure we are not inspecting anything from the generated file
		if genFile != nil && decl.Pos() >= genFile.Pos() && decl.Pos() <= genFile.End() {
			return -1, nil
		}
	}

	var valueSpecs []*EnumValueSpec
	{
		specs := slices.Filter(decl.Specs, func(v ast.Spec, idx int) bool {
			_, ok := v.(*ast.ValueSpec)
			return ok
		})
		valueSpecs = slices.Map(specs, func(v ast.Spec, idx int) *EnumValueSpec {
			return &EnumValueSpec{Node: v.(*ast.ValueSpec)}
		})
	}
	if len(valueSpecs) == 0 {
		// we only care about const blocks with value specs
		return -1, nil
	}

	// first enum type found will be assumed relevant type for an entire block
	var relevantEnumType *EnumType
	{ // find relevant enum type and assert it is the first declaration within the block
		typedConstIdx := slices.FindIndex(valueSpecs, func(vs *EnumValueSpec, _ int) bool {
			vstyp := vs.GetTypeVia(typesInfo)
			typeIdx := slices.FindIndex(enumTypes, func(et *EnumType, _ int) bool {
				return types.IdenticalIgnoreTags(vstyp, et.GetTypeVia(typesInfo))
			})
			if typeIdx == -1 {
				return false
			}
			relevantEnumType = enumTypes[typeIdx]
			return true
		})
		if typedConstIdx == -1 {
			// no enum was found, that's ok
			return -1, nil
		}
		// enums must always be defined as a pure and coherent blocks
		// we need to assert that here
		if typedConstIdx != 0 {
			// therefore we expect the first declaration to define the block's relevant enum type
			// hence this case needs to be reported and the block will not be validated any further
			return valueSpecs[typedConstIdx].Node.Pos(), errors.New("enum constants must be defined in a block of their own")
		}
		if relevantEnumType.ConstBlock != nil {
			// another block was assigned already
			// but we want all of them to be assigned in one common block
			return decl.Pos(), errors.New("enum constants must be defined in a common block")
		}
	}

	// assign block to enum type definition
	relevantEnumType.ConstBlock = &EnumConstBlock{Specs: valueSpecs, Node: decl}
	return -1, nil
}

func getTypeFromString(typeValue string) (types.BasicKind, bool) {
	typ, ok := primitiveTypes[typeValue]
	return typ, ok
}

var (
	primitiveTypes = map[string]types.BasicKind{
		"uint":       types.Uint,
		"uint8":      types.Uint8,
		"uint16":     types.Uint16,
		"uint32":     types.Uint32,
		"uint64":     types.Uint64,
		"int":        types.Int,
		"int8":       types.Int8,
		"int16":      types.Int16,
		"int32":      types.Int32,
		"int64":      types.Int64,
		"float32":    types.Float32,
		"float64":    types.Float64,
		"complex64":  types.Complex64,
		"complex128": types.Complex128,
		"bool":       types.Bool,
		"string":     types.String,
	}
	primitiveTypesReverse = map[types.BasicKind]string{
		types.Uint:       "uint",
		types.Uint8:      "uint8",
		types.Uint16:     "uint16",
		types.Uint32:     "uint32",
		types.Uint64:     "uint64",
		types.Int:        "int",
		types.Int8:       "int8",
		types.Int16:      "int16",
		types.Int32:      "int32",
		types.Int64:      "int64",
		types.Float32:    "float32",
		types.Float64:    "float64",
		types.Complex64:  "complex64",
		types.Complex128: "complex128",
		types.Bool:       "bool",
		types.String:     "string",
	}

	typedParserFuncs = map[types.BasicKind]func(raw string) (any, error){
		types.Uint: func(raw string) (any, error) {
			if len(raw) == 0 {
				return uint(0), nil
			}
			v, err := strconv.ParseUint(raw, 10, 0)
			if err != nil {
				return uint(0), err
			}
			return uint(v), err
		},
		types.Uint8: func(raw string) (any, error) {
			if len(raw) == 0 {
				return uint8(0), nil
			}
			v, err := strconv.ParseUint(raw, 10, 8)
			if err != nil {
				return uint8(0), err
			}
			return uint8(v), err
		},
		types.Uint16: func(raw string) (any, error) {
			if len(raw) == 0 {
				return uint16(0), nil
			}
			v, err := strconv.ParseUint(raw, 10, 16)
			if err != nil {
				return uint16(0), err
			}
			return uint16(v), err
		},
		types.Uint32: func(raw string) (any, error) {
			if len(raw) == 0 {
				return uint32(0), nil
			}
			v, err := strconv.ParseUint(raw, 10, 32)
			if err != nil {
				return uint32(0), err
			}
			return uint32(v), err
		},
		types.Uint64: func(raw string) (any, error) {
			if len(raw) == 0 {
				return uint64(0), nil
			}
			v, err := strconv.ParseUint(raw, 10, 64)
			if err != nil {
				return uint64(0), err
			}
			return uint64(v), err
		},
		types.Int: func(raw string) (any, error) {
			if len(raw) == 0 {
				return int(0), nil
			}
			v, err := strconv.ParseInt(raw, 10, 0)
			if err != nil {
				return int(0), err
			}
			return int(v), err
		},
		types.Int8: func(raw string) (any, error) {
			if len(raw) == 0 {
				return int8(0), nil
			}
			v, err := strconv.ParseInt(raw, 10, 8)
			if err != nil {
				return int8(0), err
			}
			return int8(v), err
		},
		types.Int16: func(raw string) (any, error) {
			if len(raw) == 0 {
				return int16(0), nil
			}
			v, err := strconv.ParseInt(raw, 10, 16)
			if err != nil {
				return int16(0), err
			}
			return int16(v), err
		},
		types.Int32: func(raw string) (any, error) {
			if len(raw) == 0 {
				return int32(0), nil
			}
			v, err := strconv.ParseInt(raw, 10, 32)
			if err != nil {
				return int32(0), err
			}
			return int32(v), err
		},
		types.Int64: func(raw string) (any, error) {
			if len(raw) == 0 {
				return int64(0), nil
			}
			v, err := strconv.ParseInt(raw, 10, 64)
			if err != nil {
				return int64(0), err
			}
			return int64(v), err
		},
		types.Float32: func(raw string) (any, error) {
			if len(raw) == 0 {
				return float32(0), nil
			}
			v, err := strconv.ParseFloat(raw, 32)
			if err != nil {
				return float32(0), err
			}
			if math.IsNaN(v) {
				err = ErrIsNaN
			} else if math.IsInf(v, 0) {
				err = ErrIsPosInf
			} else if math.IsInf(v, 0) {
				err = ErrIsNegInf
			}
			return float32(v), err
		},
		types.Float64: func(raw string) (any, error) {
			if len(raw) == 0 {
				return float64(0), nil
			}
			v, err := strconv.ParseFloat(raw, 64)
			if err != nil {
				return float64(0), err
			}
			if math.IsNaN(v) {
				err = ErrIsNaN
			} else if math.IsInf(v, 0) {
				err = ErrIsPosInf
			} else if math.IsInf(v, 0) {
				err = ErrIsNegInf
			}
			return float64(v), err
		},
		types.Complex64: func(raw string) (any, error) {
			if len(raw) == 0 {
				return complex64(0), nil
			}
			v, err := strconv.ParseComplex(raw, 64)
			if err != nil {
				return complex64(0), err
			}
			return complex64(v), err
		},
		types.Complex128: func(raw string) (any, error) {
			if len(raw) == 0 {
				return complex128(0), nil
			}
			v, err := strconv.ParseComplex(raw, 128)
			if err != nil {
				return complex128(0), err
			}
			return complex128(v), err
		},
		types.Bool: func(raw string) (any, error) {
			if len(raw) == 0 {
				return false, nil
			}
			v, err := strconv.ParseBool(raw)
			if err != nil {
				return false, err
			}
			return v, err
		},
		types.String: func(raw string) (any, error) { return raw, nil },
	}
)

func TypeToString(t types.BasicKind) string {
	return primitiveTypesReverse[t]
}

var (
	ErrIsNaN    = errors.New("typed value is NaN")
	ErrIsPosInf = errors.New("typed value is +Inf")
	ErrIsNegInf = errors.New("typed value is -Inf")
)
