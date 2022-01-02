package pills

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
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
	type testCase struct {
		p               Pill
		from            string
		serialized      string
		wantStringerOut string
		invalid         bool
		shadowedBy      Pill
	}
	testCases := []testCase{
		{from: "", serialized: "", p: Pill(-1), invalid: true, wantStringerOut: "Pill(-1)"},
		{from: "", serialized: "", p: Pill(5), invalid: true, wantStringerOut: "Pill(4)"},
		{from: "PLACEBO", serialized: "PLACEBO", p: Pill(0), wantStringerOut: "PLACEBO"},
		{from: "ASPIRIN", serialized: "ASPIRIN", p: PillAspirin, wantStringerOut: "ASPIRIN"},
		{from: "IBUPROFEN", serialized: "IBUPROFEN", p: PillIbuprofen, wantStringerOut: "IBUPROFEN"},
		{from: "PARACETAMOL", serialized: "PARACETAMOL", p: PillParacetamol, wantStringerOut: "PARACETAMOL"},
		{from: "ACETAMINOPHEN", serialized: "ACETAMINOPHEN", p: PillAcetaminophen, wantStringerOut: "ACETAMINOPHEN", shadowedBy: PillParacetamol},
		{from: "VITAMIN-C", serialized: "VITAMIN-C", p: PillVitaminC, wantStringerOut: "VITAMIN-C"},
	}
	for idx, tC := range testCases {
		t.Run(fmt.Sprintf("Serializers (idx: %d %s)", idx, tC.p), func(t *testing.T) {
			shadow := func(tC testCase) (after testCase) {
				if tC.shadowedBy == 0 {
					return tC
				}
				return testCase{
					p:               tC.shadowedBy,
					from:            tC.from,
					serialized:      tC.shadowedBy.String(),
					wantStringerOut: tC.shadowedBy.String(),
					invalid:         tC.invalid,
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
					var p Pill
					err := p.UnmarshalBinary([]byte(tC.from))
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
					var p Pill
					jsonSerialized, err := json.Marshal(shadow.from)
					require.NoError(t, err)
					err = p.UnmarshalJSON([]byte(jsonSerialized))
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, shadow.p, p)
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
					var p Pill
					err := p.UnmarshalText([]byte(tC.from))
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
					var p Pill
					err := p.UnmarshalYAML(&yaml.Node{Kind: yaml.ScalarNode, Value: tC.serialized})
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
					values := []interface{}{tC.serialized, []byte(tC.from), stringer{tC.serialized}}
					for _, v := range values {
						var p Pill
						err := p.Scan(v)
						if tC.invalid {
							require.Error(t, err)
							return
						}
						require.NoError(t, err)
						require.Equal(t, shadow.p, p)
					}
				})
				t.Run("Scan <nil>", func(t *testing.T) {
					var value interface{} = nil
					var p Pill
					err := p.Scan(value)
					require.Error(t, err)
					require.Zero(t, p)
				})
			})
		})
	}
}

type stringer struct{ v string }

func (s stringer) String() string { return s.v }
