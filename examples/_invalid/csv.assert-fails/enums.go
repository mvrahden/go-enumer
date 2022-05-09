package invalid

//go:enum -from=source.csv
type AssertConstantsAgainstCSV uint

const (
	//go:enum assert={"0":"Some Wrong Assertion"}
	AssertionFailsCSV AssertConstantsAgainstCSV = 0
)
