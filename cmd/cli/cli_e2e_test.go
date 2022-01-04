package cli_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mvrahden/go-enumer/cmd/cli"
	"github.com/stretchr/testify/require"
)

func TestE2E(t *testing.T) {
	testcases := []struct {
		dirName   string
		typealias string
		args      []string
	}{
		{"greeting", "Greeting", nil},
	}
	for idx, tC := range testcases {
		t.Run(fmt.Sprintf("Generate (idx: %d %q/%q)", idx, tC.dirName, tC.typealias), func(t *testing.T) {
			t.Cleanup(cli.CleanUpPackage)

			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, tC.dirName+"_enumer.go")
			cli.PatchTargetFilename(t, tmpFile)

			err := cli.Execute([]string{"-typealias=" + tC.typealias, "-dir=testdata/" + tC.dirName})
			require.NoError(t, err)
			require.FileExists(t, tmpFile)

			actual, err := os.ReadFile(tmpFile)
			require.NoError(t, err)
			expected, err := os.ReadFile("testdata/" + tC.dirName + "/" + strings.ToLower(tC.typealias) + "_enumer.go")
			require.NoError(t, err)
			require.Equal(t, string(expected), string(actual))
		})
	}
}
