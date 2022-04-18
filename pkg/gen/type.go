package gen

type GoType int

const (
	GoTypeUnknown GoType = -1

	GoTypeUnsignedInteger GoType = iota + 1
	GoTypeUnsignedInteger8
	GoTypeUnsignedInteger16
	GoTypeUnsignedInteger32
	GoTypeUnsignedInteger64
)

var (
	typeMap = map[string]GoType{
		"uint":   GoTypeUnsignedInteger,
		"uint8":  GoTypeUnsignedInteger8,
		"uint16": GoTypeUnsignedInteger16,
		"uint32": GoTypeUnsignedInteger32,
		"uint64": GoTypeUnsignedInteger64,
	}
)
