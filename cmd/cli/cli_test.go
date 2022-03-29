package cli

import (
	"path/filepath"
	"testing"

	"github.com/mvrahden/go-enumer/config"
	"github.com/stretchr/testify/require"
)

func TestCli(t *testing.T) {
	t.Run("generate filename", func(t *testing.T) {
		require.Equal(t,
			"/path/to/enumer.go",
			targetFilename("/path/to", &config.Options{}))
	})
	t.Run("input validation fails", func(t *testing.T) {
		t.Cleanup(CleanUpPackage)
		testcases := []struct {
			desc string
			args []string
			msg  string
		}{
			{
				"on conflicting yaml serializers",
				[]string{"-serializers=yaml,yaml.v3"},
				"serializers \"yaml\" and \"yaml.v3\" are cannot be applied together.",
			},
		}
		for _, tC := range testcases {
			t.Run(tC.desc, func(t *testing.T) {
				t.Cleanup(CleanUpPackage)
				err := Execute(tC.args)
				require.EqualError(t, err, "invalid arguments. err: "+tC.msg)
			})
		}
	})
}

func PatchTargetFilenameFunc(t *testing.T, targetPath string) {
	targetFilename = func(file string, cfg *config.Options) string {
		return filepath.Clean(targetPath)
	}
}

func CleanUpPackage() {
	cArgs = config.Args{}
	scanPath = ""
}
