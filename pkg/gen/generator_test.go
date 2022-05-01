package gen

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mvrahden/go-enumer/about"
	"github.com/mvrahden/go-enumer/config"
)

const (
	packageBase = about.Repo
)

func TestGenerator(t *testing.T) {
	for _, tC := range []struct {
		directory   string
		description string
	}{
		{"greetings", "standard enum and enum with default value"},
		{"pills", "compatibility for various unsigned integer types and forms of assignment"},
		{"planets", "standard enum and enum with default value support `ignore-case` and `undefined`"},
		{"booking", "CSV source"},
		{"colors", "CSV source with typed additional data"},
		{"project", "a set of more realistic use cases"},
	} {
		pkg := path.Join(packageBase, "examples", tC.directory)
		testdatadir := filepath.Join("..", "..", "examples", tC.directory)

		t.Run(fmt.Sprintf("Generate for package %q with %s", tC.directory, tC.description), func(t *testing.T) {
			expected := getExpectedOutputFile(t, testdatadir)
			cfg := getConfig(t, testdatadir)

			g := NewGenerator(NewInspector(cfg), NewRenderer(cfg))
			srcs, err := g.Generate(pkg)
			require.NoError(t, err)
			require.Equal(t, expected, string(srcs))
		})
	}
}

func TestGeneratorEdgeCaseDetection(t *testing.T) {
	for _, tC := range []struct {
		directory string
		errMsg    string
		cfg       config.Options
	}{
		{directory: "noninteger",
			errMsg: "Enum type of \"NonInteger\" must be of an unsigned integer type, found \"float32\"."},
		{directory: "lowerbound",
			errMsg: "Enum type of \"LowerBound\" must be of an unsigned integer type, found \"int\"."},
		{directory: "upperbound",
			errMsg: "Enum \"UpperBound\" must start with either 0 or 1."},
		{directory: "noncontinuous",
			errMsg: "Enum \"NonContinuous\" must be a continuous sequence with linear increments of 1."},
		{directory: "noncontinuous2",
			errMsg: "Enum \"NonContinuous2\" must be a continuous sequence with linear increments of 1."},
		{directory: "docstring",
			errMsg: "Failed parsing doc-string for \"InvalidDocstring\". err: flag provided but not defined: -unsupported"},
		{directory: "csv.missing-file",
			errMsg: "Failed reading from CSV for \"MissingCSV\". err: no such file \"source.csv\""},
		{directory: "csv.empty",
			errMsg: "Failed reading from CSV for \"EmptyCSV\". err: found empty csv source"},
		{directory: "csv.invalid-header",
			errMsg: "Failed reading from CSV for \"NumericFirstCellInCSV\". err: first row must be a header row but found numeric value in first cell"},
		{directory: "csv.invalid-value",
			errMsg: "Failed reading from CSV for \"NegativeValueInCSV\". err: failed converting \"-1\" to uint64"},
		{directory: "csv.invalid-range",
			errMsg: "Enum \"InvalidRangeCSV\" must be a continuous sequence with linear increments of 1."},
	} {
		t.Run(fmt.Sprintf("Generate for package %q", tC.directory), func(t *testing.T) {
			pkg := path.Join(packageBase, "examples", "_invalid", tC.directory)

			g := NewGenerator(NewInspector(&tC.cfg), NewRenderer(&tC.cfg))
			srcs, err := g.Generate(pkg)
			require.Error(t, err)
			require.ErrorContains(t, err, tC.errMsg)
			require.Zero(t, srcs)
		})
	}
}

func getConfig(t *testing.T, testdatadir string) *config.Options {
	cfg := config.LoadFrom(filepath.Join(testdatadir, "/config.yml"))
	require.NotZero(t, cfg)
	return cfg
}

func getExpectedOutputFile(t *testing.T, testdatadir string) string {
	f, err := os.Open(filepath.Join(testdatadir, "/generated.go"))
	require.NoError(t, err)
	defer f.Close()

	buf, err := io.ReadAll(f)
	require.NoError(t, err)
	els := bytes.SplitN(buf, []byte("\n"), 2)
	require.Len(t, els, 2)
	firstLine := els[0]
	require.True(t, matchGeneratedFileRegex.Match(firstLine), "Must be a generated file!")
	firstLine = []byte(fmt.Sprintf(string(firstLine), about.ShortInfo()))
	return string(firstLine) + "\n" + string(els[1])
}
