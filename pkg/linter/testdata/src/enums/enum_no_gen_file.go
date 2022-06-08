package enums

// ValidX is a valid enum, but misses a generated file
//go:enum
type ValidX uint // want `please generate enum file`

const ValidXConst ValidX = iota
