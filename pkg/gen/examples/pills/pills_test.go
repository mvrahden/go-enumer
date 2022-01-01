package pills

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPills(t *testing.T) {
	t.Run("Value Sets", func(t *testing.T) {
		require.Equal(t,
			[]string{"PLACEBO", "ASPIRIN", "IBUPROFEN", "PARACETAMOL", "VITAMIN-C"},
			PillStrings())
		require.Equal(t,
			[]Pill{PillPlacebo, PillAspirin, PillIbuprofen, PillParacetamol, PillVitaminC},
			PillValues())
	})
	t.Run("Lookup", func(t *testing.T) {
		type testCase struct {
			enum  Pill
			upper string
			lower string
		}
		testCases := []testCase{
			{PillPlacebo, "PLACEBO", "placebo"},
			{PillAspirin, "ASPIRIN", "aspirin"},
			{PillIbuprofen, "IBUPROFEN", "ibuprofen"},
			{PillParacetamol, "PARACETAMOL", "paracetamol"},
			{PillVitaminC, "VITAMIN-C", "vitamin-c"},
		}
		for idx, tC := range testCases {
			t.Run(fmt.Sprintf("Case-sensitive lookup (idx: %d %s)", idx, tC.enum), func(t *testing.T) {
				p, ok := PillFromString(tC.upper)
				require.True(t, ok)
				require.Equal(t, tC.enum, p)
				p, ok = PillFromString(tC.lower)
				require.False(t, ok)
				require.Equal(t, Pill(0), p)
			})
			t.Run(fmt.Sprintf("Case-insensitive lookup (idx: %d %s)", idx, tC.enum), func(t *testing.T) {
				p, ok := PillFromStringIgnoreCase(tC.upper)
				require.True(t, ok)
				require.Equal(t, tC.enum, p)
				p, ok = PillFromStringIgnoreCase(tC.lower)
				require.True(t, ok)
				require.Equal(t, tC.enum, p)
			})
		}
	})
	t.Run("Serialization", func(t *testing.T) {
		type testCase struct {
			p          Pill
			serialized string
			stringer   string
			invalid    bool
			shadowedBy Pill
		}
		testCases := []testCase{
			{serialized: "", p: Pill(-1), invalid: true, stringer: "Pill(-1)"},
			{serialized: "", p: Pill(5), invalid: true, stringer: "Pill(4)"},
			{serialized: "PLACEBO", p: Pill(0), stringer: "PLACEBO"},
			{serialized: "ASPIRIN", p: PillAspirin, stringer: "ASPIRIN"},
			{serialized: "IBUPROFEN", p: PillIbuprofen, stringer: "IBUPROFEN"},
			{serialized: "PARACETAMOL", p: PillParacetamol, stringer: "PARACETAMOL"},
			{serialized: "ACETAMINOPHEN", p: PillAcetaminophen, stringer: "ACETAMINOPHEN", shadowedBy: PillParacetamol},
			{serialized: "VITAMIN-C", p: PillVitaminC, stringer: "VITAMIN-C"},
		}
		for _, tC := range testCases {
			shadow := func(tC testCase) (after testCase) {
				if tC.shadowedBy == 0 {
					return tC
				}
				return testCase{
					p:          tC.shadowedBy,
					serialized: tC.shadowedBy.String(),
					stringer:   tC.shadowedBy.String(),
					invalid:    tC.invalid,
				}
			}(tC)
			t.Run("binary", func(t *testing.T) {
				t.Run("MarhsalBinary", func(t *testing.T) {
					actual, err := tC.p.MarshalBinary()
					if tC.invalid {
						require.Errorf(t, err, "got: %s", actual)
						return
					}
					require.NoError(t, err)
					require.Equal(t, shadow.serialized, string(actual))
				})
				t.Run("UnmarshalBinary", func(t *testing.T) {
					p := tC.p
					err := p.UnmarshalBinary([]byte(tC.serialized))
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, shadow.p, p)
				})
			})
			t.Run("json", func(t *testing.T) {
				t.Run("MarhsalJSON", func(t *testing.T) {
					actual, err := tC.p.MarshalJSON()
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					expected, err := json.Marshal(shadow.serialized)
					require.NoError(t, err)
					require.Equal(t, expected, actual)
				})
				t.Run("UnmarshalJSON", func(t *testing.T) {
					cp := tC.p
					jsonSerialized, err := json.Marshal(shadow.serialized)
					require.NoError(t, err)
					err = cp.UnmarshalJSON([]byte(jsonSerialized))
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, shadow.p, cp)
				})
			})
			t.Run("text", func(t *testing.T) {
				t.Run("MarhsalText", func(t *testing.T) {
					actual, err := tC.p.MarshalText()
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, shadow.serialized, string(actual))
				})
				t.Run("UnmarshalText", func(t *testing.T) {
					p := tC.p
					err := p.UnmarshalText([]byte(tC.serialized))
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, tC.p, p)
				})
			})
			t.Run("yaml", func(t *testing.T) {
				t.Run("MarhsalYAML", func(t *testing.T) {
					j, err := tC.p.MarshalYAML()
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, shadow.serialized, j)
				})
				t.Run("UnmarshalYAML", func(t *testing.T) {
					p := tC.p
					err := p.UnmarshalYAML(func(i interface{}) error {
						return json.Unmarshal([]byte("\""+tC.serialized+"\""), i)
					})
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, shadow.p, p)
				})
			})
			t.Run("sql", func(t *testing.T) {
				t.Run("Value", func(t *testing.T) {
					j, err := tC.p.Value()
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, shadow.serialized, j)
				})
				t.Run("Scan", func(t *testing.T) {
					values := []interface{}{tC.serialized, []byte(tC.serialized), stringer{tC.serialized}}
					for _, v := range values {
						p := tC.p
						err := p.Scan(v)
						if tC.invalid {
							require.Error(t, err)
							return
						}
						require.NoError(t, err)
						require.Equal(t, shadow.p, p)
					}
				})
			})
		}
	})
}

type stringer struct{ v string }

func (s stringer) String() string { return s.v }
