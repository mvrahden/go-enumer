package invalid

//go:enum -from=source.csv
type ConstOutOfRangeCSV uint

const (
	NoSuchValue ConstOutOfRangeCSV = 2
)
