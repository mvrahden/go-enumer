package noninteger

//go:enumer -from=source.csv
type MissingCSV uint

const (
	MissingCSVA MissingCSV = 0
	MissingCSVB
)
