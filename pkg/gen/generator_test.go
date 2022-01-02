package gen

import (
	"bytes"
	"embed"
	"fmt"
	"io"
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

//go:embed examples
var examples embed.FS

func TestGenerator(t *testing.T) {
	for _, target := range []string{
		"pills",
		"greetings",
	} {
		pkg := path.Join(packageBase, "/pkg/gen/examples/", target)
		testdatadir := filepath.Join("examples/", target)

		t.Run(fmt.Sprintf("Generate for package %q", target), func(t *testing.T) {
			expected := getExpectedOutputFile(t, testdatadir)
			cfg := getConfig(t, testdatadir)

			g := NewGenerator(NewInspector(cfg), NewRenderer(cfg))
			srcs, err := g.Generate(pkg)
			require.NoError(t, err)
			require.Equal(t, expected, string(srcs))
		})
	}
}

func TestEdgeCaseDetection(t *testing.T) {
	target := "pills"

	pkg := path.Join(packageBase, "/pkg/gen/examples/", target)

	for _, tC := range []struct {
		errMsg string
		cfg    config.Options
	}{
		{cfg: config.Options{TypeAliasName: "PillNotIntegerType"},
			errMsg: "Invalid enum set: Enum type must be an integer-like type, found \"float32\"."},
		{cfg: config.Options{TypeAliasName: "PillViolatesLowerBound"},
			errMsg: "Invalid enum set: Enums need to start with either 0 or 1."},
		{cfg: config.Options{TypeAliasName: "PillViolatesUpperBound"},
			errMsg: "Invalid enum set: Enums need to start with either 0 or 1."},
		{cfg: config.Options{TypeAliasName: "PillNotContinuous"},
			errMsg: "Invalid enum set: Enums must be a continuous sequence with linear increments of 1."},
	} {

		t.Run(fmt.Sprintf("Generate for package %q", target), func(t *testing.T) {
			g := NewGenerator(NewInspector(&tC.cfg), NewRenderer(&tC.cfg))
			srcs, err := g.Generate(pkg)
			require.Error(t, err)
			require.Contains(t, err.Error(), tC.errMsg)
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
	f, err := examples.Open(filepath.Join(testdatadir, "/generated.go"))
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
