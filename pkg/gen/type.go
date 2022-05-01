package gen

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

type GoType uint

const (
	GoTypeUnknown GoType = iota
	GoTypeUnsignedInteger
	GoTypeUnsignedInteger8
	GoTypeUnsignedInteger16
	GoTypeUnsignedInteger32
	GoTypeUnsignedInteger64
	GoTypeSignedInteger
	GoTypeSignedInteger8
	GoTypeSignedInteger16
	GoTypeSignedInteger32
	GoTypeSignedInteger64
	GoTypeFloat32
	GoTypeFloat64
	GoTypeComplex64
	GoTypeComplex128
	GoTypeBool
	GoTypeString
)

var (
	validEnumTypesMap = map[string]GoType{
		"uint":   GoTypeUnsignedInteger,
		"uint8":  GoTypeUnsignedInteger8,
		"uint16": GoTypeUnsignedInteger16,
		"uint32": GoTypeUnsignedInteger32,
		"uint64": GoTypeUnsignedInteger64,
	}

	primitiveTypes = map[string]GoType{
		"uint":       GoTypeUnsignedInteger,
		"uint8":      GoTypeUnsignedInteger8,
		"uint16":     GoTypeUnsignedInteger16,
		"uint32":     GoTypeUnsignedInteger32,
		"uint64":     GoTypeUnsignedInteger64,
		"int":        GoTypeSignedInteger,
		"int8":       GoTypeSignedInteger8,
		"int16":      GoTypeSignedInteger16,
		"int32":      GoTypeSignedInteger32,
		"int64":      GoTypeSignedInteger64,
		"float32":    GoTypeFloat32,
		"float64":    GoTypeFloat64,
		"complex64":  GoTypeComplex64,
		"complex128": GoTypeComplex128,
		"bool":       GoTypeBool,
		"string":     GoTypeString,
	}

	primitiveTypesReverse = map[GoType]string{
		GoTypeUnsignedInteger:   "uint",
		GoTypeUnsignedInteger8:  "uint8",
		GoTypeUnsignedInteger16: "uint16",
		GoTypeUnsignedInteger32: "uint32",
		GoTypeUnsignedInteger64: "uint64",
		GoTypeSignedInteger:     "int",
		GoTypeSignedInteger8:    "int8",
		GoTypeSignedInteger16:   "int16",
		GoTypeSignedInteger32:   "int32",
		GoTypeSignedInteger64:   "int64",
		GoTypeFloat32:           "float32",
		GoTypeFloat64:           "float64",
		GoTypeComplex64:         "complex64",
		GoTypeComplex128:        "complex128",
		GoTypeBool:              "bool",
		GoTypeString:            "string",
	}
)

var (
	ErrIsNaN    = errors.New("type is NaN")
	ErrIsPosInf = errors.New("type is +Inf")
	ErrIsNegInf = errors.New("type is -Inf")
)

func getParserFuncFor(typ GoType) func(raw string) (any, error) {
	switch typ {
	case GoTypeUnsignedInteger:
		return func(raw string) (any, error) {
			if len(raw) == 0 {
				return uint(0), nil
			}
			v, err := strconv.ParseUint(raw, 10, 0)
			if err != nil {
				return uint(0), err
			}
			return uint(v), err
		}
	case GoTypeUnsignedInteger8:
		return func(raw string) (any, error) {
			if len(raw) == 0 {
				return uint8(0), nil
			}
			v, err := strconv.ParseUint(raw, 10, 8)
			if err != nil {
				return uint8(0), err
			}
			return uint8(v), err
		}
	case GoTypeUnsignedInteger16:
		return func(raw string) (any, error) {
			if len(raw) == 0 {
				return uint16(0), nil
			}
			v, err := strconv.ParseUint(raw, 10, 16)
			if err != nil {
				return uint16(0), err
			}
			return uint16(v), err
		}
	case GoTypeUnsignedInteger32:
		return func(raw string) (any, error) {
			if len(raw) == 0 {
				return uint32(0), nil
			}
			v, err := strconv.ParseUint(raw, 10, 32)
			if err != nil {
				return uint32(0), err
			}
			return uint32(v), err
		}
	case GoTypeUnsignedInteger64:
		return func(raw string) (any, error) {
			if len(raw) == 0 {
				return uint64(0), nil
			}
			v, err := strconv.ParseUint(raw, 10, 64)
			if err != nil {
				return uint64(0), err
			}
			return uint64(v), err
		}
	case GoTypeSignedInteger:
		return func(raw string) (any, error) {
			if len(raw) == 0 {
				return int(0), nil
			}
			v, err := strconv.ParseInt(raw, 10, 0)
			if err != nil {
				return int(0), err
			}
			return int(v), err
		}
	case GoTypeSignedInteger8:
		return func(raw string) (any, error) {
			if len(raw) == 0 {
				return int8(0), nil
			}
			v, err := strconv.ParseInt(raw, 10, 8)
			if err != nil {
				return int8(0), err
			}
			return int8(v), err
		}
	case GoTypeSignedInteger16:
		return func(raw string) (any, error) {
			if len(raw) == 0 {
				return int16(0), nil
			}
			v, err := strconv.ParseInt(raw, 10, 16)
			if err != nil {
				return int16(0), err
			}
			return int16(v), err
		}
	case GoTypeSignedInteger32:
		return func(raw string) (any, error) {
			if len(raw) == 0 {
				return int32(0), nil
			}
			v, err := strconv.ParseInt(raw, 10, 32)
			if err != nil {
				return int32(0), err
			}
			return int32(v), err
		}
	case GoTypeSignedInteger64:
		return func(raw string) (any, error) {
			if len(raw) == 0 {
				return int64(0), nil
			}
			v, err := strconv.ParseInt(raw, 10, 64)
			if err != nil {
				return int64(0), err
			}
			return int64(v), err
		}
	case GoTypeFloat32:
		return func(raw string) (any, error) {
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
		}
	case GoTypeFloat64:
		return func(raw string) (any, error) {
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
		}
	case GoTypeComplex64:
		return func(raw string) (any, error) {
			if len(raw) == 0 {
				return complex64(0), nil
			}
			v, err := strconv.ParseComplex(raw, 64)
			if err != nil {
				return complex64(0), err
			}
			return complex64(v), err
		}
	case GoTypeComplex128:
		return func(raw string) (any, error) {
			if len(raw) == 0 {
				return complex128(0), nil
			}
			v, err := strconv.ParseComplex(raw, 128)
			if err != nil {
				return complex128(0), err
			}
			return complex128(v), err
		}
	case GoTypeBool:
		return func(raw string) (any, error) {
			if len(raw) == 0 {
				return false, nil
			}
			v, err := strconv.ParseBool(raw)
			if err != nil {
				return false, err
			}
			return v, err
		}
	case GoTypeString:
		fallthrough
	default:
		return func(raw string) (any, error) {
			return raw, nil
		}
	}
}

func getTypeFromString(typeValue string) GoType {
	typ, ok := primitiveTypes[typeValue]
	if !ok {
		return GoTypeUnknown
	}
	return typ
}

func (t GoType) String() (typeValue string) {
	v, ok := primitiveTypesReverse[t]
	if !ok {
		return "string"
	}
	return v
}

func (t GoType) ZeroValueString() (zeroValue string) {
	switch t {
	case GoTypeUnsignedInteger, GoTypeUnsignedInteger8, GoTypeUnsignedInteger16, GoTypeUnsignedInteger32, GoTypeUnsignedInteger64,
		GoTypeSignedInteger, GoTypeSignedInteger8, GoTypeSignedInteger16, GoTypeSignedInteger32, GoTypeSignedInteger64,
		GoTypeFloat32, GoTypeFloat64,
		GoTypeComplex64, GoTypeComplex128:
		return "0"
	case GoTypeBool:
		return "false"
	case GoTypeString, GoTypeUnknown:
		fallthrough
	default:
		return "\"\""
	}
}

func (t GoType) ToSource(v any) (src string) {
	switch t {
	case GoTypeString, GoTypeUnknown:
		return fmt.Sprintf("%q", v)
	default:
		return fmt.Sprintf("%s", v)
	}
}
