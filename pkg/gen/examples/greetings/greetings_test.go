package greetings

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGreetings(t *testing.T) {
	t.Run("Value Sets", func(t *testing.T) {
		require.Equal(t,
			[]string{"Россия", "中國", "日本", "한국", "ČeskáRepublika", "𝜋"},
			GreetingStrings())
		require.Equal(t,
			[]Greeting{GreetingРоссия, Greeting中國, Greeting日本, Greeting한국, GreetingČeskáRepublika, Greeting𝜋},
			GreetingValues())
		t.Run("Ent Interface", func(t *testing.T) {
			require.Equal(t,
				[]string{"Россия", "中國", "日本", "한국", "ČeskáRepublika", "𝜋"},
				Greeting(0).Values())
		})
	})
	t.Run("Lookup", func(t *testing.T) {
		type testCase struct {
			enum  Greeting
			upper string
			lower string
		}
		testCases := []testCase{
			{GreetingUndefined, "", ""},
			{GreetingРоссия, "Россия", "россия"},
			{Greeting中國, "中國", "中國"},
			{Greeting日本, "日本", "日本"},
			{Greeting한국, "한국", "한국"},
			{GreetingČeskáRepublika, "ČeskáRepublika", "ČeskáRepublika"},
			{Greeting𝜋, "𝜋", "𝜋"},
		}
		for idx, tC := range testCases {
			t.Run(fmt.Sprintf("Case-sensitive lookup (idx: %d %s)", idx, tC.enum), func(t *testing.T) {
				p, ok := GreetingFromString(tC.upper)
				require.True(t, ok)
				require.Equal(t, tC.enum, p)
				p, ok = GreetingFromString(tC.lower)
				if tC.lower == tC.upper {
					require.True(t, ok)
					require.Equal(t, tC.enum, p)
				} else {
					require.False(t, ok)
					require.Equal(t, Greeting(0), p)
				}
			})
			t.Run(fmt.Sprintf("Case-insensitive lookup (idx: %d %s)", idx, tC.enum), func(t *testing.T) {
				p, ok := GreetingFromStringIgnoreCase(tC.upper)
				require.True(t, ok)
				require.Equal(t, tC.enum, p)
				p, ok = GreetingFromStringIgnoreCase(tC.lower)
				require.True(t, ok)
				require.Equal(t, tC.enum, p)
			})
		}
	})
	testCases := []struct {
		from       string
		g          Greeting
		serialized string
		invalid    bool
		stringer   string
	}{
		{from: "", serialized: "", g: Greeting(0), invalid: false, stringer: ""},
		{from: "Greeting(7)", serialized: "Greeting(7)", g: Greeting(7), invalid: true, stringer: "Greeting(7)"},
		{from: "Россия", serialized: "Россия", g: GreetingРоссия, stringer: "россия"},
		{from: "中國", serialized: "中國", g: Greeting中國, stringer: "中國"},
		{from: "日本", serialized: "日本", g: Greeting日本, stringer: "日本"},
		{from: "한국", serialized: "한국", g: Greeting한국, stringer: "한국"},
		{from: "ČeskáRepublika", serialized: "ČeskáRepublika", g: GreetingČeskáRepublika, stringer: "ČeskáRepublika"},
		{from: "𝜋", serialized: "𝜋", g: Greeting𝜋, stringer: "𝜋"},
	}
	for idx, tC := range testCases {
		t.Run(fmt.Sprintf("Serializers (idx: %d %s)", idx, tC.g), func(t *testing.T) {
			t.Run("binary", func(t *testing.T) {
				t.Run("MarhsalBinary", func(t *testing.T) {
					j, err := tC.g.MarshalBinary()
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, tC.serialized, string(j))
				})
				t.Run("UnmarshalBinary", func(t *testing.T) {
					var g Greeting
					err := g.UnmarshalBinary([]byte(tC.from))
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, tC.g, g)
				})
			})
			t.Run("json", func(t *testing.T) {
				t.Run("MarhsalJSON", func(t *testing.T) {
					jsonSerialized, err := json.Marshal(tC.serialized)
					require.NoError(t, err)
					actual, err := tC.g.MarshalJSON()
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, jsonSerialized, actual)
				})
				t.Run("UnmarshalJSON", func(t *testing.T) {
					var g Greeting
					err := g.UnmarshalJSON([]byte("\"" + tC.from + "\""))
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, tC.g, g)
				})
			})
			t.Run("text", func(t *testing.T) {
				t.Run("MarhsalText", func(t *testing.T) {
					j, err := tC.g.MarshalText()
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, tC.serialized, string(j))
				})
				t.Run("UnmarshalText", func(t *testing.T) {
					var g Greeting
					err := g.UnmarshalText([]byte(tC.from))
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, tC.g, g)
				})
			})
			t.Run("yaml", func(t *testing.T) {
				t.Run("MarhsalYAML", func(t *testing.T) {
					j, err := tC.g.MarshalYAML()
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, tC.serialized, j)
				})
				t.Run("UnmarshalYAML", func(t *testing.T) {
					var g Greeting
					err := g.UnmarshalYAML(func(i interface{}) error {
						return json.Unmarshal([]byte("\""+tC.serialized+"\""), i)
					})
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, tC.g, g)
				})
			})
			t.Run("sql", func(t *testing.T) {
				t.Run("Value", func(t *testing.T) {
					j, err := tC.g.Value()
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, tC.serialized, j)
				})
				t.Run("Scan", func(t *testing.T) {
					values := []interface{}{tC.serialized, []byte(tC.from), stringer{tC.serialized}}
					for _, v := range values {
						var g Greeting
						err := g.Scan(v)
						if tC.invalid {
							require.Error(t, err)
							return
						}
						require.NoError(t, err)
						require.Equal(t, tC.g, g)
					}
				})
				t.Run("Scan <nil>", func(t *testing.T) {
					var value interface{} = nil
					var g Greeting
					err := g.Scan(value)
					require.NoError(t, err) // we have set "undefined"
					require.Zero(t, g)
				})
			})
		})
	}
}

func TestGreetingsWithDefault(t *testing.T) {
	t.Run("Value Sets", func(t *testing.T) {
		require.Equal(t,
			[]string{"World", "Россия", "中國", "日本", "한국", "ČeskáRepublika", "𝜋"},
			GreetingWithDefaultStrings())
		require.Equal(t,
			[]GreetingWithDefault{GreetingWithDefaultWorld, GreetingWithDefaultРоссия, GreetingWithDefault中國, GreetingWithDefault日本, GreetingWithDefault한국, GreetingWithDefaultČeskáRepublika, GreetingWithDefault𝜋},
			GreetingWithDefaultValues())
		t.Run("Ent Interface", func(t *testing.T) {
			require.Equal(t,
				[]string{"World", "Россия", "中國", "日本", "한국", "ČeskáRepublika", "𝜋"},
				GreetingWithDefault(0).Values())
		})
	})
	t.Run("Lookup", func(t *testing.T) {
		type testCase struct {
			enum  GreetingWithDefault
			upper string
			lower string
		}
		testCases := []testCase{
			{GreetingWithDefaultWorld, "", ""}, // default value
			{GreetingWithDefaultWorld, "World", "world"},
			{GreetingWithDefaultРоссия, "Россия", "россия"},
			{GreetingWithDefault中國, "中國", "中國"},
			{GreetingWithDefault日本, "日本", "日本"},
			{GreetingWithDefault한국, "한국", "한국"},
			{GreetingWithDefaultČeskáRepublika, "ČeskáRepublika", "českárepublika"},
			{GreetingWithDefault𝜋, "𝜋", "𝜋"},
		}
		for idx, tC := range testCases {
			t.Run(fmt.Sprintf("Case-sensitive lookup (idx: %d %s)", idx, tC.enum), func(t *testing.T) {
				g, ok := GreetingWithDefaultFromString(tC.upper)
				require.True(t, ok)
				require.Equal(t, tC.enum, g)
				g, ok = GreetingWithDefaultFromString(tC.lower)
				if tC.lower == tC.upper {
					require.True(t, ok)
					require.Equal(t, tC.enum, g)
				} else {
					require.False(t, ok)
					require.Equal(t, GreetingWithDefault(0), g)
				}
			})
			t.Run(fmt.Sprintf("Case-insensitive lookup (idx: %d %s)", idx, tC.enum), func(t *testing.T) {
				g, ok := GreetingWithDefaultFromStringIgnoreCase(tC.upper)
				require.True(t, ok)
				require.Equal(t, tC.enum, g)
				g, ok = GreetingWithDefaultFromStringIgnoreCase(tC.lower)
				require.True(t, ok)
				require.Equal(t, tC.enum, g)
			})
		}
	})
	testCases := []struct {
		g          GreetingWithDefault
		from       string
		serialized string
		invalid    bool
		stringer   string
	}{
		{from: "World", serialized: "World", g: GreetingWithDefault(0), invalid: false, stringer: "World"},
		{from: "GreetingWithDefault(7)", serialized: "GreetingWithDefault(7)", g: GreetingWithDefault(7), invalid: true, stringer: "GreetingWithDefault(7)"},
		{from: "", serialized: "World", g: GreetingWithDefaultWorld, stringer: "World"}, // default
		{from: "World", serialized: "World", g: GreetingWithDefaultWorld, stringer: "World"},
		{from: "Россия", serialized: "Россия", g: GreetingWithDefaultРоссия, stringer: "россия"},
		{from: "中國", serialized: "中國", g: GreetingWithDefault中國, stringer: "中國"},
		{from: "日本", serialized: "日本", g: GreetingWithDefault日本, stringer: "日本"},
		{from: "한국", serialized: "한국", g: GreetingWithDefault한국, stringer: "한국"},
		{from: "ČeskáRepublika", serialized: "ČeskáRepublika", g: GreetingWithDefaultČeskáRepublika, stringer: "ČeskáRepublika"},
		{from: "𝜋", serialized: "𝜋", g: GreetingWithDefault𝜋, stringer: "𝜋"},
	}
	for idx, tC := range testCases {
		t.Run(fmt.Sprintf("Serializers (idx: %d %s)", idx, tC.g), func(t *testing.T) {
			t.Run("binary", func(t *testing.T) {
				t.Run("MarhsalBinary", func(t *testing.T) {
					actual, err := tC.g.MarshalBinary()
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, tC.serialized, string(actual))
				})
				t.Run("UnmarshalBinary", func(t *testing.T) {
					var g GreetingWithDefault
					err := g.UnmarshalBinary([]byte(tC.from))
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, tC.g, g)
				})
			})
			t.Run("json", func(t *testing.T) {
				t.Run("MarhsalJSON", func(t *testing.T) {
					actual, err := tC.g.MarshalJSON()
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					jsonSerialized, err := json.Marshal(tC.serialized)
					require.NoError(t, err)
					require.Equal(t, jsonSerialized, actual)
				})
				t.Run("UnmarshalJSON", func(t *testing.T) {
					var g GreetingWithDefault
					err := g.UnmarshalJSON([]byte("\"" + tC.from + "\""))
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, tC.g, g)
				})
			})
			t.Run("text", func(t *testing.T) {
				t.Run("MarhsalText", func(t *testing.T) {
					actual, err := tC.g.MarshalText()
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, tC.serialized, string(actual))
				})
				t.Run("UnmarshalText", func(t *testing.T) {
					var g GreetingWithDefault
					err := g.UnmarshalText([]byte(tC.from))
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, tC.g, g)
				})
			})
			t.Run("yaml", func(t *testing.T) {
				t.Run("MarhsalYAML", func(t *testing.T) {
					actual, err := tC.g.MarshalYAML()
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, tC.serialized, actual)
				})
				t.Run("UnmarshalYAML", func(t *testing.T) {
					var g GreetingWithDefault
					err := g.UnmarshalYAML(func(i interface{}) error {
						return json.Unmarshal([]byte("\""+tC.from+"\""), i)
					})
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, tC.g, g)
				})
			})
			t.Run("sql", func(t *testing.T) {
				t.Run("Value", func(t *testing.T) {
					actual, err := tC.g.Value()
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, tC.serialized, actual)
				})
				t.Run("Scan", func(t *testing.T) {
					values := []interface{}{tC.from, []byte(tC.from), stringer{tC.from}}
					for _, v := range values {
						var g GreetingWithDefault
						err := g.Scan(v)
						if tC.invalid {
							require.Error(t, err)
							return
						}
						require.NoError(t, err)
						require.Equal(t, tC.g, g)
					}
				})
				t.Run("Scan <nil>", func(t *testing.T) {
					var value interface{} = nil
					var g GreetingWithDefault
					err := g.Scan(value)
					require.NoError(t, err) // we have set a default value
					require.Zero(t, g)
				})
			})
		})
	}
}

type stringer struct{ v string }

func (s stringer) String() string { return s.v }
