package enums

//go:enum -from=source_a.csv
type Color1 uint

const (
	Color1Hello Color1 = 0
	Color1World Color1 = 0 // want `redundant constant`
)

//go:enum -from=source_a.csv // want `enum of same file already exists`
type Color2 uint

//go:enum -from=source_b.csv
type Color3 uint

const (
	Color3Hello Color3 = 4 // this is OK for const blocks from file based enums
)

// //go:enum -from=source_b.csv
// type Color4 uint

// const (
// 	Color4Hello Color4 = 1 // assert "white" //-want-`assertion failed`

// 	// //go:enum someCommand //-want-`unknown command \"someCommand\"`
// 	// ColorInvalidA Color = 0

// 	// //go:enum assert not-a-json-object
// 	// ColorInvalidB Color = 0 //-want-`command parameter must be a valid and meaningful JSON object`
// 	// //go:enum assert {"not-a-json-object"}
// 	// ColorInvalidC Color = 0 //-want-`command parameter must be a valid and meaningful JSON object`
// 	// //go:enum assert {}
// 	// ColorInvalidD Color = 0 //-want-`command parameter must be a valid and meaningful JSON object`
// 	// //go:enum assert "Pink"

// 	// //go:enum
// 	// //go:enum
// 	// ColorBlack Color = 0 //-want-`at most one magic comment permitted`
// )
