package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigLoading(t *testing.T) {
	configFile := filepath.Join(t.TempDir(), "config.yml")
	t.Run("Load defaults", func(t *testing.T) {
		cfg := LoadFrom("")
		require.Equal(t, &Options{
			TransformStrategy: "noop",
		}, cfg)
	})
	t.Run("Load from Config file", func(t *testing.T) {
		err := os.WriteFile(configFile, []byte(`---
transform: def
addPrefix: ghi
serializers: [jkl,mno,pqr]
support: [stu,vwx,yz]`), os.ModePerm)
		require.NoError(t, err)

		cfg := LoadFrom(configFile)
		require.Equal(t, &Options{
			TransformStrategy: "def",
			AddPrefix:         "ghi",
			Serializers:       []string{"jkl", "mno", "pqr"},
			SupportedFeatures: []string{"stu", "vwx", "yz"},
		}, cfg)
	})
	t.Run("Load with Args", func(t *testing.T) {
		t.Run("joins with defaults if value not present", func(t *testing.T) {
			args := &Args{}
			cfg := LoadWith(args)
			require.Equal(t, &Options{
				TransformStrategy: "noop",
			}, cfg)
		})
		t.Run("preserves value if value present", func(t *testing.T) {
			args := &Args{TransformStrategy: "def"}
			cfg := LoadWith(args)
			require.Equal(t, &Options{
				TransformStrategy: "def",
			}, cfg)
		})
	})
}

func TestCopy(t *testing.T) {
	t.Run("copies to a new instance with equal values", func(t *testing.T) {
		cfg := LoadFrom("")
		cfg.AddPrefix = "ABC"
		cfg.Serializers = stringList{"DEF"}
		cfg2 := cfg.Clone()
		require.NotSame(t, cfg, cfg2)
		require.Equal(t, cfg, cfg2)
	})
}

func TestStringList(t *testing.T) {
	t.Run("Contains", func(t *testing.T) {
		require.False(t, stringList{"a", "b", "c"}.Contains("v"))
		require.True(t, stringList{"a", "b", "c"}.Contains("a"))
		require.True(t, stringList{"a", "b", "c"}.Contains("b"))
		require.True(t, stringList{"a", "b", "c"}.Contains("c"))
	})
	t.Run("Unique", func(t *testing.T) {
		require.Equal(t, []string{"a", "b"}, stringList{"a", "b", "a"}.Unique())
	})
	t.Run("String", func(t *testing.T) {
		require.Equal(t, "a,b,a", stringList{"a", "b", "a"}.String())
	})
	t.Run("Set", func(t *testing.T) {
		sl := stringList{"a", "b", "a"}
		require.NoError(t, sl.Set("a,b, c, d"))
		require.Equal(t, "a,b,c,d", sl.String())
	})
}
