package utils

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"gopkg.in/yaml.v3"
)

func Must[T any](a T, v any) T {
	switch t := v.(type) {
	case bool:
		if !t {
			panic(fmt.Sprintf("invalid Must() call for %#v", a))
		}
	case error:
		if t != nil {
			panic(fmt.Sprintf("invalid Must() call for %#v. got err: %s", a, t))
		}
	case nil:
		return a
	default:
		panic("invalid use of Must(). Second arg must be `bool` or `error`")
	}
	return a
}

func AssertNotSamePointer(t *testing.T, expected, actual any) {
	expPtr := fmt.Sprintf("%p", expected)
	actPtr := fmt.Sprintf("%p", actual)
	if expPtr == actPtr {
		t.Fatalf("%T and %T point to the same address %[1]p", expected, actual)
	}
}

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
}

type TestConfig struct {
	SupportUndefined bool
	HasDefault       bool
}

func AssertMissingSerializationInterfacesFor[T any](t *testing.T, missingSerializers []string) {
	for _, serializer := range missingSerializers {
		t.Run(fmt.Sprintf("Not Implemented Serializer Interface %q", serializer), func(t *testing.T) {
			assertMissingSerializer[T](t, serializer)
			assertMissingDeserializer[T](t, serializer)
		})
	}
}

func AssertSerializationInterfacesFor[T any](t *testing.T, idx int, tC TestCase, cfg TestConfig, serializers []string) {
	t.Run(fmt.Sprintf("Serializers (idx: %d from %q)", idx, tC.From), func(t *testing.T) {
		assertSerializers[T](t, tC, cfg, serializers)
	})
	t.Run(fmt.Sprintf("Deserializers (idx: %d from %q)", idx, tC.From), func(t *testing.T) {
		assertDeserializers[T](t, tC, cfg, serializers)
	})
}

func assertSerializers[T any](t *testing.T, tC TestCase, cfg TestConfig, serializers []string) {
	for _, serializer := range serializers {
		t.Run(fmt.Sprintf("serialize %q", serializer), func(t *testing.T) {
			assertSerializer[T](t, tC, cfg, serializer)
		})
	}
}

func assertSerializer[T any](t *testing.T, tC TestCase, cfg TestConfig, serializer string) {
	switch serializer {
	case "binary":
		t.Run("MarhsalBinary", func(t *testing.T) {
			enum := tC.Enum.(interface {
				MarshalBinary() (data []byte, err error)
			})
			j, err := enum.MarshalBinary()
			if tC.Expected.IsInvalid && isDefault(cfg.HasDefault, tC.Enum) {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tC.Expected.AsSerialized, string(j))
		})

	case "bson":
		t.Run("MarshalBSONValue", func(t *testing.T) {
			enum := tC.Enum.(interface {
				MarshalBSONValue() (_ bsontype.Type, data []byte, err error)
			})
			typ, actual, err := enum.MarshalBSONValue()
			if tC.Expected.IsInvalid && isDefault(cfg.HasDefault, tC.Enum) {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if cfg.SupportUndefined && !cfg.HasDefault {
				// If expected is Zero Value
				// we need Nullability
				if isZero(tC.Enum) {
					require.Equal(t, bsontype.Undefined, typ)
					require.Equal(t, []byte(nil), actual)
					return
				}
			}
			require.Equal(t, bsontype.String, typ)
			require.NotNil(t, actual)
		})

	case "graphql":
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
			if tC.Expected.IsInvalid && isDefault(cfg.HasDefault, tC.Enum) {
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
			actual, err := enum.Value()
			if tC.Expected.IsInvalid && isDefault(cfg.HasDefault, tC.Enum) {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			if cfg.SupportUndefined && !cfg.HasDefault {
				// If expected is Zero Value
				// we need Nullability
				if isZero(tC.Enum) {
					require.Equal(t, nil, actual)
					return
				}
			}
			require.Equal(t, tC.Expected.AsSerialized, actual)
		})

	case "text":
		t.Run("MarhsalText", func(t *testing.T) {
			enum := tC.Enum.(interface {
				MarshalText() (text []byte, err error)
			})
			actual, err := enum.MarshalText()
			if tC.Expected.IsInvalid && isDefault(cfg.HasDefault, tC.Enum) {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tC.Expected.AsSerialized, string(actual))
		})

	case "yaml", "yaml.v3":
		t.Run("MarhsalYAML", func(t *testing.T) {
			enum := tC.Enum.(interface {
				MarshalYAML() (any, error)
			})
			j, err := enum.MarshalYAML()
			if tC.Expected.IsInvalid && isDefault(cfg.HasDefault, tC.Enum) {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tC.Expected.AsSerialized, j)
		})

	default:
		require.FailNow(t, "unsupported serializer %q", serializer)
	}
}

func ToPointer[T any](v T) *T {
	return &v
}

// `zeroValuer` helps to constructor non-"nil" instances of pointer type T.
// It is a workaround for Limitations of Generics when it comes to mutable
// interface functions with pointer receiver.
func zeroValuer[T any]() *T {
	var v T
	return &v
}

func isZero(v any) bool {
	val := reflect.ValueOf(v)
	return val.IsValid() && val.Elem().IsZero()
}

func isDefault(hasDefault bool, v any) bool { return !(hasDefault && isZero(v)) }

func assertDeserializers[T any](t *testing.T, tC TestCase, cfg TestConfig, deserializers []string) {
	for _, deserializer := range deserializers {
		t.Run(fmt.Sprintf("deserialize %q", deserializer), func(t *testing.T) {
			assertDeserializer[T](t, tC, cfg, deserializer)
		})
	}
}

// `zeroValuer` helps to constructor non-"nil" instances of pointer type T.
func assertDeserializer[T any](t *testing.T, tC TestCase, cfg TestConfig, deserializer string) {
	switch deserializer {
	case "binary":
		t.Run("UnmarshalBinary", func(t *testing.T) {
			enum := zeroValuer[T]()
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

	case "bson":
		t.Run("UnmarshalBSONValue", func(t *testing.T) {
			enum := zeroValuer[T]()
			typ, buf, err := bson.MarshalValue(tC.From)
			require.NoError(t, err)
			require.NotNil(t, buf)
			require.Equal(t, bsontype.String, typ)

			err = (any)(enum).(interface {
				UnmarshalBSONValue(t bsontype.Type, data []byte) error
			}).UnmarshalBSONValue(bsontype.String, buf)
			if tC.Expected.IsInvalid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tC.Enum, enum)
		})

	case "graphql":
		t.Run("UnmarshalGQL", func(t *testing.T) {
			values := []any{tC.From, []byte(tC.From), stringer{tC.From}}
			for _, v := range values {
				enum := zeroValuer[T]()
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
			enum := zeroValuer[T]()
			err := (any)(enum).(interface {
				UnmarshalGQL(value any) error
			}).UnmarshalGQL(v)
			require.Equal(t, zeroValuer[T](), enum)
			if !cfg.SupportUndefined {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})

	case "json":
		t.Run("UnmarshalJSON", func(t *testing.T) {
			enum := zeroValuer[T]()
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
				enum := zeroValuer[T]()
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
			enum := zeroValuer[T]()
			err := (any)(enum).(interface {
				Scan(src any) error
			}).Scan(v)
			require.Equal(t, zeroValuer[T](), enum)
			if !cfg.SupportUndefined {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})

	case "text":
		t.Run("UnmarshalText", func(t *testing.T) {
			enum := zeroValuer[T]()
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
			enum := zeroValuer[T]()
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

	case "yaml.v3":
		t.Run("UnmarshalYAML", func(t *testing.T) {
			enum := zeroValuer[T]()
			err := (any)(enum).(interface {
				UnmarshalYAML(n *yaml.Node) error
			}).UnmarshalYAML(&yaml.Node{Tag: "!!str", Value: tC.From})
			if tC.Expected.IsInvalid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tC.Enum, enum)
		})
	default:
		require.FailNow(t, "unsupported deserializer %q", deserializer)
	}
}

func assertMissingSerializer[T any](t *testing.T, serializer string) {
	var ok bool
	switch serializer {
	case "binary":
		t.Run("MarhsalBinary", func(t *testing.T) {
			var enum T
			_, ok = (any)(enum).(interface {
				MarshalBinary() (data []byte, err error)
			})
		})

	case "graphql":
		t.Run("MarhsalGQL", func(t *testing.T) {
			var enum T
			_, ok = (any)(enum).(interface {
				MarshalGQL(w io.Writer)
			})
		})

	case "json":
		t.Run("MarhsalJSON", func(t *testing.T) {
			var enum T
			_, ok = (any)(enum).(interface {
				MarshalJSON() ([]byte, error)
			})
		})

	case "sql":
		t.Run("Value (SQL)", func(t *testing.T) {
			var enum T
			_, ok = (any)(enum).(interface {
				Value() (driver.Value, error)
			})
		})

	case "text":
		t.Run("MarhsalText", func(t *testing.T) {
			var enum T
			_, ok = (any)(enum).(interface {
				MarshalText() (text []byte, err error)
			})
		})

	case "yaml", "yaml.v3":
		t.Run("MarhsalYAML", func(t *testing.T) {
			var enum T
			_, ok = (any)(enum).(interface {
				MarshalYAML() (any, error)
			})
		})

	default:
		require.FailNow(t, "unsupported serializer %q", serializer)
	}
	require.Falsef(t, ok, "Expected to NOT implement interface for %q", serializer)
}

func assertMissingDeserializer[T any](t *testing.T, deserializer string) {
	var ok bool
	switch deserializer {
	case "binary":
		t.Run("UnmarshalBinary", func(t *testing.T) {
			_, ok = (any)(zeroValuer[T]()).(interface {
				UnmarshalBinary(data []byte) error
			})
		})

	case "graphql":
		t.Run("UnmarshalGQL", func(t *testing.T) {
			_, ok = (any)(zeroValuer[T]()).(interface {
				UnmarshalGQL(value any) error
			})
		})
		t.Run("UnmarshalGQL from <nil>", func(t *testing.T) {
			_, ok = (any)(zeroValuer[T]()).(interface {
				UnmarshalGQL(value any) error
			})
		})

	case "json":
		t.Run("UnmarshalJSON", func(t *testing.T) {
			_, ok = (any)(zeroValuer[T]()).(interface {
				UnmarshalJSON([]byte) error
			})
		})

	case "sql":
		t.Run("Scan (SQL)", func(t *testing.T) {
			_, ok = (any)(zeroValuer[T]()).(interface {
				Scan(src any) error
			})
		})

	case "text":
		t.Run("UnmarshalText", func(t *testing.T) {
			_, ok = (any)(zeroValuer[T]()).(interface {
				UnmarshalText(text []byte) error
			})
		})

	case "yaml":
		t.Run("UnmarshalYAML", func(t *testing.T) {
			_, ok = (any)(zeroValuer[T]()).(interface {
				UnmarshalYAML(unmarshal func(any) error) error
			})
		})

	case "yaml.v3":
		t.Run("UnmarshalYAML", func(t *testing.T) {
			_, ok = (any)(zeroValuer[T]()).(interface {
				UnmarshalYAML(n *yaml.Node) error
			})
		})
	default:
		require.FailNow(t, "unsupported deserializer %q", deserializer)
	}
	require.Falsef(t, ok, "Expected to NOT implement interface for %q", deserializer)
}
