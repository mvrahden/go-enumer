package cli_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/mvrahden/go-enumer/cmd/cli"
	"github.com/stretchr/testify/require"
)

func TestE2E_Cli(t *testing.T) {
	cli.PatchDeleteOldGeneratedFileFunc(t)
	t.Cleanup(func() {
		cli.ResetDeleteOldGeneratedFileFunc(t)
	})

	testcases := []struct {
		desc        string
		dirName     string
		args        []string
		outFilename string
	}{
		{"no args",
			"greeting", nil, "gen.golden"},
		{"standard serializers",
			"greeting", []string{"-serializers=binary,graphql,json,sql,text,yaml"}, "gen.serializers.golden"},
		{"standard serializers - different order, same result",
			"greeting", []string{"-serializers=sql,graphql,json,yaml,binary,text"}, "gen.serializers.golden"},
		{"serializers and yaml.v3",
			"greeting", []string{"-serializers=binary,graphql,json,sql,text,yaml.v3"}, "gen.serializers.yaml_v3.golden"},
		{"standard serializers - deserialize with ignore case",
			"greeting", []string{"-serializers=sql,graphql,json,yaml,binary,text", "-support=ignore-case"}, "gen.serializers.ignore-case.golden"},
		{"standard output and ent interface",
			"greeting", []string{"-support=ent"}, "gen.ent.golden"},
	}
	for idx, tC := range testcases {
		t.Run(fmt.Sprintf("Generate (idx: %d %q)", idx, tC.desc), func(t *testing.T) {

			tmpDir := t.TempDir()
			cli.PatchTargetFilenameFunc(t, tmpDir)
			tmpFile := filepath.Join(tmpDir, tC.outFilename+".go")

			args := []string{"-dir=testdata/" + tC.dirName, "-out=" + tC.outFilename}
			args = append(args, tC.args...)
			err := cli.Execute(args)
			require.NoError(t, err)
			require.FileExists(t, tmpFile)

			actual, err := os.ReadFile(tmpFile)
			require.NoError(t, err)
			expected, err := os.ReadFile(filepath.Join("testdata", tC.dirName, tC.outFilename))
			require.NoError(t, err)
			require.Equal(t, string(expected), string(actual))
		})
	}
}

func TestE2E_DeleteOldGeneratedFile(t *testing.T) {
	t.Run("delete generated file from temp directory with various files", func(t *testing.T) {
		tmpDir := t.TempDir()

		// unrelated file - this file stays
		err := os.WriteFile(filepath.Join(tmpDir, "testFile"), []byte("I'm not generated"), os.ModePerm)
		require.NoError(t, err)

		{ // write some files
			buf, err := ioutil.ReadFile(filepath.Join("testdata", "greeting", "gen.golden"))
			require.NoError(t, err)

			// this file has a comment but is not a Go file - it stays
			err = os.Mkdir(filepath.Join(tmpDir, "keepMe_Dir"), os.ModeDir)
			require.NoError(t, err)

			err = os.WriteFile(filepath.Join(tmpDir, "keepMe"), buf, os.ModePerm)
			require.NoError(t, err)
			// no/incomplete gen comment match - it stays
			err = os.WriteFile(filepath.Join(tmpDir, "keepMe_0.go"), buf[2:], os.ModePerm)
			require.NoError(t, err)
			// < 78 bytes - it stays
			err = os.WriteFile(filepath.Join(tmpDir, "keepMe_1.go"), buf[:77], os.ModePerm)
			require.NoError(t, err)
			// > 78 but < 85 bytes - it stays
			err = os.WriteFile(filepath.Join(tmpDir, "keepMe_2.go"), buf[:81], os.ModePerm)
			require.NoError(t, err)
			// this is our TARGET (marked with x to ensure its read as last entry)
			err = os.WriteFile(filepath.Join(tmpDir, "x_deleteMe.go"), buf[:100], os.ModePerm)
			require.NoError(t, err)
		}

		args := []string{"-dir=" + tmpDir}
		err = cli.Execute(args)
		require.ErrorContains(t, err, "no enums detected")

		require.NoFileExists(t, filepath.Join(tmpDir, "deleteMe.go"))

		require.DirExists(t, filepath.Join(tmpDir, "keepMe_Dir"))
		for _, filename := range []string{"keepMe", "keepMe_0.go", "keepMe_1.go", "keepMe_2.go"} {
			require.FileExists(t, filepath.Join(tmpDir, filename))
		}
	})
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
			"on no enums in given directory (wrong path)", []string{"-dir=testdata/nothing-here"}, "no such directory",
		},
		{
			"on invalid enum sequence", []string{"-dir=testdata/error_cases/non_continuous_sequence"}, "\"InvalidNonContinuousGreeting\" type specification is invalid. err: enum const block must not contain skipped rows",
		},
	}
	for _, tC := range testcases {
		t.Run(tC.desc, func(t *testing.T) {
			err := cli.Execute(tC.args)
			require.ErrorContains(t, err, tC.msg)
		})
	}
}
