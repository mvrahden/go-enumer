package enums

//go:enum -from=source_a.csv
type Color1 uint // want `\"Color1Hello\" is a redundant constant`

const (
	Color1Hello Color1 = 0
	Color1World Color1 = 0
)

//go:enum -from=source_a.csv // want `enum of same file already exists`
type Color2 uint

//go:enum -from=source_b.csv
type Color3 uint // want `\"Color3Hello\" exceeds spec range \[0,0\]`

const (
	Color3Hello Color3 = 4
)

//go:enum -from=source_c.csv
type Color4 uint // want `\"Color4Redundant1\" is a redundant constant`

const (
	Color4Redundant1 Color4 = 2
	Color4Redundant2 Color4 = 2
)

//go:enum -from=source_d.csv
type Color5 uint // want `\"Color5Green\" fails on assertion (reason: assertion failed)`

const (
	Color5White Color5 = 1 // assert "White"
	Color5Green Color5 = 3 // assert "Red"
)
