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
	"github.com/mvrahden/go-enumer/pkg/common"
)

const (
	packageBase = about.Repo
)

func TestGeneratorExamples(t *testing.T) {
	for _, tC := range []struct {
		directory   string
		description string
	}{
		{"greetings", "standard enum and enum with default value"},
		{"animals", "standard enums with some transformations"},
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
			errMsg: "\"NonInteger\" type specification is invalid. err: enum types must be of any unsigned integer type"},
		{directory: "lowerbound",
			errMsg: "\"LowerBound\" type specification is invalid. err: enum types must be of any unsigned integer type"},
		{directory: "upperbound",
			errMsg: "\"UpperBound\" type specification is invalid. err: enum spec sequences must start with either 0 or 1"},
		{directory: "noncontinuous",
			errMsg: "\"NonContinuous\" type specification is invalid. err: enum spec sequences must increment at most by one"},
		{directory: "noncontinuous2",
			errMsg: "\"NonContinuous2\" type specification is invalid. err: enum const block must not contain skipped rows"},
		{directory: "noncontinuous3",
			errMsg: "\"NonContinuous3\" type specification is invalid. err: enum spec sequences must increment at most by one"},
		{directory: "rowed",
			errMsg: "\"Rowed\" type specification is invalid. err: enum const block must not contain rowed declarations"},
		{directory: "block-unrelated-types",
			errMsg: "\"Unrelated\" type specification is invalid. err: enum const block must not contain unrelated type declarations"},
		{directory: "docstring",
			errMsg: "\"InvalidDocstring\" type specification is invalid. err: unknown option \"unsupported\""},
		{directory: "csv.no-path-traversal",
			errMsg: "\"ForbiddenPathTraversalCSV\" type specification is invalid. err: source path cannot contain path traversals"},
		{directory: "csv.no-path-traversal-2",
			errMsg: "\"ForbiddenPathTraversalCSV\" type specification is invalid. err: source path cannot contain path traversals"},
		{directory: "csv.no-relative-path-prefix",
			errMsg: "\"NoRelativePathPrefixCSV\" type specification is invalid. err: source path cannot start with \"./\" or \"/\""},
		{directory: "csv.empty",
			errMsg: "\"EmptyCSV\" type specification is invalid. err: found empty csv source"},
		{directory: "csv.invalid-header",
			errMsg: "\"NumericFirstCellInCSV\" type specification is invalid. err: header cannot contain numeric values"},
		{directory: "csv.invalid-value",
			errMsg: "\"NegativeValueInCSV\" type specification is invalid. err: failed converting \"-1\" to uint64"},
		{directory: "csv.missing-file",
			errMsg: "\"MissingCSV\" type specification is invalid. err: no such file"},
		{directory: "csv.missing-value-row",
			errMsg: "\"NumericFirstCellInCSV\" type specification is invalid. err: csv source must contain at least one value row"},
		{directory: "csv.range-start",
			errMsg: "\"InvalidRangeCSV\" type specification is invalid. err: enum sequences must start with either 0 or 1"},
		{directory: "csv.range-noncontinuous",
			errMsg: "\"InvalidRangeCSV\" type specification is invalid. err: enum spec sequences must increment at most by one"},
		{directory: "csv.range-noncontinuous-2",
			errMsg: "\"InvalidRangeCSV\" type specification is invalid. err: enum spec sequences must increment at most by one"},
		{directory: "csv.const-out-of-range",
			errMsg: "\"ConstOutOfRangeCSV\" type specification is invalid. err: \"NoSuchValue\" exceeds spec range [1,1]"},
		{directory: "csv.const-out-of-range-2",
			errMsg: "\"ConstOutOfRangeCSV\" type specification is invalid. err: \"NoSuchValue\" exceeds spec range [0,1]"},
		{directory: "csv.assertion-failed",
			errMsg: "\"AssertionFailedCSV\" type specification is invalid. err: \"NotAnApple\" fails on assertion (reason: assertion failed)"},
		{directory: "csv.assertion-failed-2",
			errMsg: "\"AssertionFailedCSV\" type specification is invalid. err: \"NotAnApple\" fails on assertion (reason: missing terminating quote in assertion)"},
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
	require.True(t, common.GEN_ENUMER_FILE.Match(firstLine), "Must be a generated file!")
	return string(buf)
}
