package invalid

//go:enum
type NonInteger float32

const (
	NonIntegerA NonInteger = iota
	NonIntegerB
)
