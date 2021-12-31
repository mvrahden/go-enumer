package gen

type GoType int

const (
	GoTypeUnknown GoType = -1

	GoTypeSignedInteger GoType = iota + 1
	GoTypeSignedInteger8
	GoTypeSignedInteger16
	GoTypeSignedInteger32
	GoTypeSignedInteger64
	GoTypeUnsignedInteger
	GoTypeUnsignedInteger8
	GoTypeUnsignedInteger16
	GoTypeUnsignedInteger32
	GoTypeUnsignedInteger64
)

var (
	typeMap = map[string]GoType{
		"int":    GoTypeSignedInteger,
		"int8":   GoTypeSignedInteger8,
		"int16":  GoTypeSignedInteger16,
		"int32":  GoTypeSignedInteger32,
		"int64":  GoTypeSignedInteger64,
		"uint":   GoTypeUnsignedInteger,
		"uint8":  GoTypeUnsignedInteger8,
		"uint16": GoTypeUnsignedInteger16,
		"uint32": GoTypeUnsignedInteger32,
		"uint64": GoTypeUnsignedInteger64,
	}
	typeMapReverse = map[GoType]string{
		GoTypeSignedInteger:     "int",
		GoTypeSignedInteger8:    "int8",
		GoTypeSignedInteger16:   "int16",
		GoTypeSignedInteger32:   "int32",
		GoTypeSignedInteger64:   "int64",
		GoTypeUnsignedInteger:   "uint",
		GoTypeUnsignedInteger8:  "uint8",
		GoTypeUnsignedInteger16: "uint16",
		GoTypeUnsignedInteger32: "uint32",
		GoTypeUnsignedInteger64: "uint64",
	}
)
