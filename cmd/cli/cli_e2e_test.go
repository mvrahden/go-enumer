package cli_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/mvrahden/go-enumer/cmd/cli"
	"github.com/stretchr/testify/require"
)

func TestE2E_Cli(t *testing.T) {
	testcases := []struct {
		desc        string
		dirName     string
		args        []string
		outFilename string
	}{
		{"no args",
			"greeting", nil, "gen.golden"},
		{"standard serializers",
			"greeting", []string{"-serializers=binary,gql,json,sql,text,yaml"}, "gen.serializers.golden"},
		{"standard serializers - different order, same result",
			"greeting", []string{"-serializers=sql,gql,json,yaml,binary,text"}, "gen.serializers.golden"},
		{"serializers and yaml.v3",
			"greeting", []string{"-serializers=binary,gql,json,sql,text,yaml.v3"}, "gen.serializers.yaml_v3.golden"},
		{"standard serializers - deserialize with ignore case",
			"greeting", []string{"-serializers=sql,gql,json,yaml,binary,text", "-support=ignore-case"}, "gen.serializers.ignore-case.golden"},
		{"standard output and ent interface",
			"greeting", []string{"-support=ent"}, "gen.ent.golden"},
	}
	for idx, tC := range testcases {
		t.Run(fmt.Sprintf("Generate (idx: %d %q)", idx, tC.desc), func(t *testing.T) {
			t.Cleanup(cli.CleanUpPackage)

			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, tC.outFilename)
			cli.PatchTargetFilenameFunc(t, tmpFile)

			defaultArgs := []string{"-dir=testdata/" + tC.dirName}
			err := cli.Execute(append(defaultArgs, tC.args...))
			require.NoError(t, err)
			require.FileExists(t, tmpFile)

			actual, err := os.ReadFile(tmpFile)
			require.NoError(t, err)
			expected, err := os.ReadFile("testdata/" + tC.dirName + "/" + tC.outFilename)
			require.NoError(t, err)
			require.Equal(t, string(expected), string(actual))
		})
	}
}

func TestE2E_Errors(t *testing.T) {
	testcases := []struct {
		desc string
		args []string
		msg  string
	}{
		{
			"on no enums in CWD", nil, "no enums detected.",
		},
		{
			"on no enums in given directory (wrong path)", []string{"-dir=testdata/nothing-here"}, "no enums detected.",
		},
		{
			"on invalid enum sequence", []string{"-dir=testdata/error_cases/non_continuous_sequence"}, "Enum \"InvalidNonContinuousGreeting\" must be a continuous sequence with linear increments of 1.",
		},
	}
	for _, tC := range testcases {
		t.Run(tC.desc, func(t *testing.T) {
			t.Cleanup(cli.CleanUpPackage)
			err := cli.Execute(tC.args)
			require.EqualError(t, err, "failed generating code. err: "+tC.msg)
		})
	}
}
