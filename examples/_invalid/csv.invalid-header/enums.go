package invalid

//go:enum -from=source.csv
type NumericFirstCellInCSV uint

const (
	NumericFirstCellInCSVA NumericFirstCellInCSV = 0
	NumericFirstCellInCSVB
)
