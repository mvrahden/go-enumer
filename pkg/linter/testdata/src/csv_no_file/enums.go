package enums

//go:enum -from=source.xyz // want `unsupported file extension`
type EnumType1 uint

//go:enum -from=source.csv // want `no such file`
type EnumType2 uint
