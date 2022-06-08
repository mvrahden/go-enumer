package enums

//go:enum
type D string // want `enum types must be of any unsigned integer type`

//go:enum
type E rune // want `enum types must be of any unsigned integer type`

//go:enum
type F int // want `enum types must be of any unsigned integer type`

//go:enum
type G []byte // want `enum types must be of any unsigned integer type`

//go:enum
type H []any // want `enum types must be of any unsigned integer type`

//go:enum
type I map[any]any // want `enum types must be of any unsigned integer type`

//go:enum
type J func() // want `enum types must be of any unsigned integer type`

//go:enum
type K chan any // want `enum types must be of any unsigned integer type`

//go:enum
type L struct{} // want `enum types must be of any unsigned integer type`

//go:enum
type M float32 // want `enum types must be of any unsigned integer type`

//go:enum
type N float64 // want `enum types must be of any unsigned integer type`

//go:enum
type O complex64 // want `enum types must be of any unsigned integer type`

//go:enum
type P complex128 // want `enum types must be of any unsigned integer type`

//go:enum
type R *uint // want `enum types must be of any unsigned integer type`
