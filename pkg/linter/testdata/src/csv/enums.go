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
type Color4 uint // want `\"Color4Hello\" fails on assertion \(reason: assertion failed\)`

const (
	Color4Hello Color4 = 1 // assert "white"
)
