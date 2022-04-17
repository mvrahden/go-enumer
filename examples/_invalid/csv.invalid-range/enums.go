package noninteger

//go:enum -from=source.csv
type InvalidRangeCSV uint

const (
	InvalidRangeCSVA InvalidRangeCSV = 0
	InvalidRangeCSVB
)
