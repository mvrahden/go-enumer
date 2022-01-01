package greetings

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGreetings(t *testing.T) {
	t.Run("Value Sets", func(t *testing.T) {
		require.Equal(t,
			[]string{"россия", "中國", "日本", "한국", "ČeskáRepublika", "𝜋"},
			GreetingStrings())
		require.Equal(t,
			[]Greeting{Greetingроссия, Greeting中國, Greeting日本, Greeting한국, GreetingČeskáRepublika, Greeting𝜋},
			GreetingValues())
	})
	t.Run("Serialization", func(t *testing.T) {
		testCases := []struct {
			g          Greeting
			serialized string
			invalid    bool
			stringer   string
		}{
			{serialized: "", g: Greeting(0), invalid: true, stringer: "Greeting(0)"},
			{serialized: "", g: Greeting(7), invalid: true, stringer: "Greeting(7)"},
			{serialized: "россия", g: Greetingроссия, stringer: "россия"},
			{serialized: "中國", g: Greeting中國, stringer: "中國"},
			{serialized: "日本", g: Greeting日本, stringer: "日本"},
			{serialized: "한국", g: Greeting한국, stringer: "한국"},
			{serialized: "ČeskáRepublika", g: GreetingČeskáRepublika, stringer: "ČeskáRepublika"},
			{serialized: "𝜋", g: Greeting𝜋, stringer: "𝜋"},
		}
		for _, tC := range testCases {
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
					p := tC.g
					err := p.UnmarshalBinary([]byte(tC.serialized))
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, tC.g, p)
				})
			})
			t.Run("json", func(t *testing.T) {
				jsonSerialized, err := json.Marshal(tC.serialized)
				require.NoError(t, err)
				t.Run("MarhsalJSON", func(t *testing.T) {
					actual, err := tC.g.MarshalJSON()
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, jsonSerialized, actual)
				})
				t.Run("UnmarshalJSON", func(t *testing.T) {
					p := tC.g
					err := p.UnmarshalJSON(jsonSerialized)
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, tC.g, p)
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
					p := tC.g
					err := p.UnmarshalText([]byte(tC.serialized))
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, tC.g, p)
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
					p := tC.g
					err := p.UnmarshalYAML(func(i interface{}) error {
						return json.Unmarshal([]byte("\""+tC.serialized+"\""), i)
					})
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, tC.g, p)
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
					p := tC.g
					err := p.Scan(tC.serialized)
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, tC.g, p)
				})
			})
		}
	})
}
