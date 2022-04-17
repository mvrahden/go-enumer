package noninteger

//go:enum -from=source.csv
type MissingCSV uint

const (
	MissingCSVA MissingCSV = 0
	MissingCSVB
)
