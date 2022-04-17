package noninteger

//go:enum -unsupported=value
type InvalidDocstring uint

const (
	InvalidDocstringA InvalidDocstring = 0
	InvalidDocstringB
)
