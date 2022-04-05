package utils

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

type stringer struct{ v string }

func (s stringer) String() string { return s.v }

type TestCase struct {
	From     string
	Enum     any
	Expected Expected
}

type Expected struct {
	AsSerialized string
	IsInvalid    bool
	// _IsNillable  bool
}

type TestConfig struct {
	SupportUndefined bool
}

func AssertSerializers[T any](t *testing.T, tC TestCase, assert string) {
	switch assert {
	case "binary":
		t.Run("MarhsalBinary", func(t *testing.T) {
			enum := tC.Enum.(interface {
				MarshalBinary() (data []byte, err error)
			})
			j, err := enum.MarshalBinary()
			if tC.Expected.IsInvalid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tC.Expected.AsSerialized, string(j))
		})

	case "gql":
		t.Run("MarhsalGQL", func(t *testing.T) {
			enum := tC.Enum.(interface {
				MarshalGQL(w io.Writer)
			})
			var buf bytes.Buffer
			enum.MarshalGQL(&buf)
			require.Equal(t, fmt.Sprintf("%q", tC.Expected.AsSerialized), buf.String())
		})

	case "json":
		t.Run("MarhsalJSON", func(t *testing.T) {
			jsonSerialized, err := json.Marshal(tC.Expected.AsSerialized)
			require.NoError(t, err)
			enum := tC.Enum.(interface {
				MarshalJSON() ([]byte, error)
			})
			actual, err := enum.MarshalJSON()
			if tC.Expected.IsInvalid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, jsonSerialized, actual)
		})

	case "sql":
		t.Run("Value (SQL)", func(t *testing.T) {
			enum := tC.Enum.(interface {
				Value() (driver.Value, error)
			})
			j, err := enum.Value()
			if tC.Expected.IsInvalid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tC.Expected.AsSerialized, j)
		})

	case "text":
		t.Run("MarhsalText", func(t *testing.T) {
			enum := tC.Enum.(interface {
				MarshalText() (text []byte, err error)
			})
			j, err := enum.MarshalText()
			if tC.Expected.IsInvalid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tC.Expected.AsSerialized, string(j))
		})

	case "yaml":
		t.Run("MarhsalYAML", func(t *testing.T) {
			enum := tC.Enum.(interface {
				MarshalYAML() (any, error)
			})
			j, err := enum.MarshalYAML()
			if tC.Expected.IsInvalid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tC.Expected.AsSerialized, j)
		})
	default:
		require.FailNow(t, "invalid input")
	}
}

func ToPointer[T any](v T) *T {
	return &v
}

// `zeroValuer` helps to constructor non-"nil" instances of pointer type T.
// It is a workaround for Limitations of Generics when it comes to mutable
// interface functions with pointer receiver.
func ZeroValuer[T any]() *T {
	var v T
	return &v
}

// `zeroValuer` helps to constructor non-"nil" instances of pointer type T.
func AssertDeserializers[T any](t *testing.T, tC TestCase, cfg TestConfig, assert string, zeroValuer func() T) {
	switch assert {
	case "binary":
		t.Run("UnmarshalBinary", func(t *testing.T) {
			enum := zeroValuer()
			err := (any)(enum).(interface {
				UnmarshalBinary(data []byte) error
			}).UnmarshalBinary([]byte(tC.From))
			if tC.Expected.IsInvalid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tC.Enum, enum)
		})

	case "gql":
		t.Run("UnmarshalGQL", func(t *testing.T) {
			values := []any{tC.From, []byte(tC.From), stringer{tC.From}}
			for _, v := range values {
				enum := zeroValuer()
				err := (any)(enum).(interface {
					UnmarshalGQL(value any) error
				}).UnmarshalGQL(v)
				if tC.Expected.IsInvalid {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
				require.Equal(t, tC.Enum, enum)
			}
		})
		t.Run("UnmarshalGQL from <nil>", func(t *testing.T) {
			var v any = nil
			enum := zeroValuer()
			err := (any)(enum).(interface {
				UnmarshalGQL(value any) error
			}).UnmarshalGQL(v)
			require.Equal(t, zeroValuer(), enum)
			if !cfg.SupportUndefined {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})

	case "json":
		t.Run("UnmarshalJSON", func(t *testing.T) {
			enum := zeroValuer()
			err := (any)(enum).(interface {
				UnmarshalJSON([]byte) error
			}).UnmarshalJSON([]byte("\"" + tC.From + "\""))
			if tC.Expected.IsInvalid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tC.Enum, enum)
		})

	case "sql":
		t.Run("Scan (SQL)", func(t *testing.T) {
			values := []any{tC.From, []byte(tC.From), stringer{tC.From}}
			for _, v := range values {
				enum := zeroValuer()
				err := (any)(enum).(interface {
					Scan(src any) error
				}).Scan(v)
				if tC.Expected.IsInvalid {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
				require.Equal(t, tC.Enum, enum)
			}
		})
		t.Run("Scan from <nil> (SQL)", func(t *testing.T) {
			var v any = nil
			enum := zeroValuer()
			err := (any)(enum).(interface {
				Scan(src any) error
			}).Scan(v)
			require.Equal(t, zeroValuer(), enum)
			if !cfg.SupportUndefined {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})

	case "text":
		t.Run("UnmarshalText", func(t *testing.T) {
			enum := zeroValuer()
			err := (any)(enum).(interface {
				UnmarshalText(text []byte) error
			}).UnmarshalText([]byte(tC.From))
			if tC.Expected.IsInvalid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tC.Enum, enum)
		})

	case "yaml":
		t.Run("UnmarshalYAML", func(t *testing.T) {
			enum := zeroValuer()
			err := (any)(enum).(interface {
				UnmarshalYAML(unmarshal func(any) error) error
			}).UnmarshalYAML(func(i any) error {
				return json.Unmarshal([]byte("\""+tC.From+"\""), i)
			})
			if tC.Expected.IsInvalid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tC.Enum, enum)
		})
	default:
		require.FailNow(t, "invalid deserializer %q")
	}
}
