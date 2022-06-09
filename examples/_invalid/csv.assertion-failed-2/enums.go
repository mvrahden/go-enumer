package invalid

//go:enum -from=source.csv
type AssertionFailedCSV uint

const (
	NotAnApple AssertionFailedCSV = 2 // assert "Apple // Missing terminating Quote
)
