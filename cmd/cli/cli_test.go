package cli

import (
	"testing"

	"github.com/mvrahden/go-enumer/config"
	"github.com/stretchr/testify/require"
)

func TestCli(t *testing.T) {
	t.Run("generate filename", func(t *testing.T) {
		require.Equal(t,
			"/path/to/types_enumer.go",
			targetFilename("/path/to", "types_enumer", &config.Options{}))
	})
	t.Run("input validation fails", func(t *testing.T) {
		testcases := []struct {
			desc string
			args []string
			msg  string
		}{
			{
				"on empty output file name",
				[]string{"-out="},
				"output file name cannot be empty",
			},
			{
				"on spaces in output file name",
				[]string{"-out=\"hello dude\""},
				"output file name contains spaces",
			},
			{
				"on spaces in output file name",
				[]string{"-out=\"\""},
				"output file name contains forbidden characters",
			},
			{
				"on conflicting yaml serializers",
				[]string{"-serializers=yaml,yaml.v3"},
				"serializers \"yaml\" and \"yaml.v3\" cannot be applied together",
			},
		}
		for _, tC := range testcases {
			t.Run(tC.desc, func(t *testing.T) {
				err := Execute(tC.args)
				require.EqualError(t, err, "invalid arguments. err: "+tC.msg)
			})
		}
	})
}

var deleteOldGeneratedFileFuncBackup = findAndDeleteOldGeneratedFile

func PatchDeleteOldGeneratedFileFunc(t *testing.T) {
	findAndDeleteOldGeneratedFile = func(dir string) error {
		return nil
	}
}

func ResetDeleteOldGeneratedFileFunc(t *testing.T) {
	findAndDeleteOldGeneratedFile = deleteOldGeneratedFileFuncBackup
}

var targetFilenameFuncBackup = targetFilename

func PatchTargetFilenameFunc(t *testing.T, targetDirectory string) {
	targetFilename = func(dir, file string, cfg *config.Options) string {
		return targetFilenameFuncBackup(targetDirectory, file, cfg)
	}
}
