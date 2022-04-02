package noninteger

//go:enumer -from=source.csv
type NegativeValueInCSV uint

const (
	NegativeValueInCSVA NegativeValueInCSV = 0
	NegativeValueInCSVB
)
