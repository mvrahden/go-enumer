package noninteger

//go:enumer -from=source.csv
type EmptyCSV uint

const (
	EmptyCSVA EmptyCSV = 0
	EmptyCSVB
)
