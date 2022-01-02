package about

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTags(t *testing.T) {
	t.Run("Git Info", func(t *testing.T) {
		require.NotEmpty(t, GitFormat)
		require.NotEqual(t, "local", GitFormat)
		require.NotEmpty(t, GitCommit)
		require.NotEqual(t, "local", GitCommit)
		require.NotEmpty(t, GitBranch)
		require.NotEqual(t, "local", GitBranch)
		require.NotEmpty(t, GitTag)
		require.NotEqual(t, "local", GitTag)
	})
	t.Run("Go Info", func(t *testing.T) {
		require.NotEmpty(t, GoVersion)
		require.Contains(t, []string{"windows", "darwin", "linux"}, GoOS)
		require.Contains(t, []string{"amd64"}, GoArch)
	})
	t.Run("Application Info", func(t *testing.T) {
		require.NotEmpty(t, Repo)
		require.NotEmpty(t, Application)
	})
	t.Run("Application Info Banner", func(t *testing.T) {
		require.Contains(t, ShortInfo(), Application)
		require.Contains(t, ShortInfo(), Repo)
		require.Contains(t, LongInfo(), Application)
		require.Contains(t, LongInfo(), Repo)
		require.Contains(t, LongInfo(), GitTag)
		require.Contains(t, LongInfo(), GoVersion)
		require.Contains(t, LongInfo(), GoOS)
		require.Contains(t, LongInfo(), GoArch)
	})
}
