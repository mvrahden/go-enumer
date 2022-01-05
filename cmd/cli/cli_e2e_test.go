package cli_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/mvrahden/go-enumer/cmd/cli"
	"github.com/stretchr/testify/require"
)

func TestE2E(t *testing.T) {
	testcases := []struct {
		desc        string
		dirName     string
		typealias   string
		args        []string
		outFilename string
	}{
		{"no args",
			"greeting", "Greeting", nil, "gen.golden"},
		{"standard serializers",
			"greeting", "Greeting", []string{"-serializers=binary,gql,json,sql,text,yaml"}, "gen.serializers.golden"},
		{"standard serializers - different order, same result",
			"greeting", "Greeting", []string{"-serializers=sql,gql,json,yaml,binary,text"}, "gen.serializers.golden"},
		{"serializers and yaml.v3",
			"greeting", "Greeting", []string{"-serializers=binary,gql,json,sql,text,yaml.v3"}, "gen.serializers.yaml_v3.golden"},
		{"standard serializers - deserialize with ignore case",
			"greeting", "Greeting", []string{"-serializers=sql,gql,json,yaml,binary,text", "-support=ignore-case"}, "gen.serializers.ignore-case.golden"},
		{"standard output and ent interface",
			"greeting", "Greeting", []string{"-support=ent"}, "gen.ent.golden"},
	}
	for idx, tC := range testcases {
		t.Run(fmt.Sprintf("Generate (idx: %d %q)", idx, tC.desc), func(t *testing.T) {
			t.Cleanup(cli.CleanUpPackage)

			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, tC.outFilename)
			cli.PatchTargetFilename(t, tmpFile)

			defaultArgs := []string{"-typealias=" + tC.typealias, "-dir=testdata/" + tC.dirName}
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
