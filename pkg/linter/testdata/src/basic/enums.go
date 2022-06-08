package enums

// ValidA
//go:enum -support=undefined
type ValidA uint

const (
	ValidAHello ValidA = iota + 1
	ValidAWorld
)

// ValidB
//go:enum
type ValidB uint

const (
	ValidBHello ValidB = iota + 1
	ValidBWorld
)

// ValidC
//go:enum
type ValidC uint

const (
	ValidCHello ValidC = iota
	ValidCWorld
	ValidCFoo = ValidCWorld
	ValidCBar
	ValidCBar2
)

// ValidD
//go:enum
type ValidD uint

const (
	ValidDHello ValidD = iota
	ValidDWorld
	ValidDFoo        = ValidDWorld // not nice (alt value in middle), but valid
	ValidDBar ValidD = iota - 1
	ValidDBar2
)

// ValidE
//go:enum
type ValidE uint

const (
	ValidEHello ValidE = iota
	ValidEWorld
	ValidEFoo        = ValidE(2) // should this be valid?
	ValidEBar ValidE = iota
	ValidEBar2
)

//go:enum -from=/abc/source // want `unsupported file extension`
type ValidF0 uint

//go:enum -from=/abc/source.csv // want `source path cannot start with \"./\" or \"/\"`
type ValidF1 uint

//go:enum -from=./abc/source.csv // want `source path cannot start with \"./\" or \"/\"`
type ValidF2 uint

//go:enum -from=../abc/source.csv // want `source path cannot contain path traversals`
type ValidG1 uint

//go:enum -from=abc/../../source.csv // want `source path cannot contain path traversals`
type ValidG2 uint
