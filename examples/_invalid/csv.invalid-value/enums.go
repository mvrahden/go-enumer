package noninteger

//go:enum -from=source.csv
type NegativeValueInCSV uint

const (
	NegativeValueInCSVA NegativeValueInCSV = 0
	NegativeValueInCSVB
)
