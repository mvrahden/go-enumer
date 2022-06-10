package enums

//go:enum // want `magic comment must be last row of doc string for enum type`
// InvalidA
type InvalidA uint

// InvalidB
//go:enum
//go:enum
type InvalidB uint // want `at most one magic comment permitted per enum type`

// InvalidC
//go:enum
type InvalidC uint // want `enum types require a const block or a file source`

const (
	anotherConst  string   = "hello"
	InvalidCHello InvalidC = iota // want `enum constants must be defined in a block of their own`
)

// InvalidD
//go:enum
type InvalidD uint // want `enum const block must not contain rowed declarations`

const (
	InvalidDHello              InvalidD = iota
	InvalidDWorld, InvalidDFoo InvalidD = 0, 1
	InvalidDBar, InvalidDBaz            = 0, 1
)

// InvalidE
//go:enum
type InvalidE uint // want `enum const block must not contain rowed declarations`

const (
	InvalidEWorld, InvalidEFoo InvalidE = 0, 1
	InvalidEBar, InvalidEBaz            = 0, 1
)

// InvalidF
//go:enum
type InvalidF uint // want `enum const block must not contain unrelated type declarations`

const (
	InvalidFHello InvalidF = iota
	InvalidFWorld
	InvalidFFoo = 2
	InvalidFBar = 3
)

const (
	InvalidG0Hello InvalidG0 = 0
)

// InvalidG0
//go:enum
type InvalidG0 uint // want `enum const block must be defined after its type definition`

// InvalidG1
//go:enum
type InvalidG1 uint

const (
	InvalidG1Foo InvalidG1 = iota + 1
)

const ( // want `enum constants must be defined in a common block`
	InvalidG1Bar InvalidG1 = 0
)

// InvalidG2 has const in other file
//go:enum
type InvalidG2 uint // want `enum const block must be in same file as their type definition`

// InvalidH
//go:enum
type InvalidH uint // want `enum spec sequences must start with either 0 or 1`

const (
	InvalidHHello InvalidH = iota + 2
)

// InvalidI
//go:enum
type InvalidI uint // want `enum spec sequences must be ordered`

const (
	InvalidIWorld InvalidI = 0
	InvalidIHello InvalidI = 2
	InvalidIFoo   InvalidI = 1
)

// InvalidJ
//go:enum
type InvalidJ uint // want `enum spec sequences must start with either 0 or 1`

const (
	InvalidJHello InvalidJ = iota + 2
	InvalidJWorld
)

// InvalidK
//go:enum
type InvalidK uint // want `enum spec sequences must increment at most by one`

const (
	InvalidKHello InvalidK = iota
	InvalidKWorld
	InvalidKFoo InvalidK = 3
)

// InvalidL
//go:enum arg1 arg2 // want `unknown args \[arg1 arg2\]`
type InvalidL uint

// InvalidM
//go:enum -someOption=abc // want `unknown option \"someOption\"`
type InvalidM uint
