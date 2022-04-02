package noninteger

//go:enumer -unsupported=value
type InvalidDocstring uint

const (
	InvalidDocstringA InvalidDocstring = 0
	InvalidDocstringB
)
