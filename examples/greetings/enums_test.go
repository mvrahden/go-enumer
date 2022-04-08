package greetings

import (
	"fmt"
	"testing"

	"github.com/mvrahden/go-enumer/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestEnums(t *testing.T) {
	t.Run("Greeting", func(t *testing.T) {
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
		t.Run("Serialization", func(t *testing.T) {
			cfg := utils.TestConfig{SupportUndefined: true}
			toPtr := utils.ToPointer[Greeting]
			testCases := []utils.TestCase{
				{From: "", Enum: toPtr(0), Expected: utils.Expected{AsSerialized: ""}},
				{From: "UNKNOWN", Enum: toPtr(7), Expected: utils.Expected{AsSerialized: "Greeting(7)", IsInvalid: true}},
				{From: "Россия", Enum: toPtr(GreetingРоссия), Expected: utils.Expected{AsSerialized: "Россия"}},
				{From: "中國", Enum: toPtr(Greeting中國), Expected: utils.Expected{AsSerialized: "中國"}},
				{From: "日本", Enum: toPtr(Greeting日本), Expected: utils.Expected{AsSerialized: "日本"}},
				{From: "한국", Enum: toPtr(Greeting한국), Expected: utils.Expected{AsSerialized: "한국"}},
				{From: "ČeskáRepublika", Enum: toPtr(GreetingČeskáRepublika), Expected: utils.Expected{AsSerialized: "ČeskáRepublika"}},
				{From: "𝜋", Enum: toPtr(Greeting𝜋), Expected: utils.Expected{AsSerialized: "𝜋"}},
			}
			for idx, tC := range testCases {
				serializers := []string{"binary", "gql", "json", "sql", "text", "yaml"}
				utils.AssertSerializationInterfacesFor[Greeting](t, idx, tC, cfg, serializers)
			}
		})
	})
	t.Run("GreetingWithDefault", func(t *testing.T) {
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
		t.Run("Serialization", func(t *testing.T) {
			cfg := utils.TestConfig{SupportUndefined: true}
			toPtr := utils.ToPointer[GreetingWithDefault]
			testCases := []utils.TestCase{
				{From: "", Enum: toPtr(0), Expected: utils.Expected{AsSerialized: "World"}},
				{From: "World", Enum: toPtr(0), Expected: utils.Expected{AsSerialized: "World"}},
				{From: "UNKNOWN", Enum: toPtr(7), Expected: utils.Expected{AsSerialized: "GreetingWithDefault(7)", IsInvalid: true}},
				{From: "Россия", Enum: toPtr(GreetingWithDefaultРоссия), Expected: utils.Expected{AsSerialized: "Россия"}},
				{From: "中國", Enum: toPtr(GreetingWithDefault中國), Expected: utils.Expected{AsSerialized: "中國"}},
				{From: "日本", Enum: toPtr(GreetingWithDefault日本), Expected: utils.Expected{AsSerialized: "日本"}},
				{From: "한국", Enum: toPtr(GreetingWithDefault한국), Expected: utils.Expected{AsSerialized: "한국"}},
				{From: "ČeskáRepublika", Enum: toPtr(GreetingWithDefaultČeskáRepublika), Expected: utils.Expected{AsSerialized: "ČeskáRepublika"}},
				{From: "𝜋", Enum: toPtr(GreetingWithDefault𝜋), Expected: utils.Expected{AsSerialized: "𝜋"}},
			}
			for idx, tC := range testCases {
				serializers := []string{"binary", "gql", "json", "sql", "text", "yaml"}
				utils.AssertSerializationInterfacesFor[GreetingWithDefault](t, idx, tC, cfg, serializers)
			}
		})
	})
}
