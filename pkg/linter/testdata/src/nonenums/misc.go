package nonenums

type A string

type B int

type C uint

const (
	// these constants are a prove that non-enum types are skipped
	// they violate enum rules, as they have:
	// - a start of > 1
	// - a non-linear increase
	C1 C = iota + 2
	C2
	C3
	C4 = 10
)
