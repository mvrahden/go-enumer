package invalid

//go:enum
type NonContinuous3 uint

const (
	NonContinuous3A NonContinuous3 = iota << 2
	NonContinuous3B
)
