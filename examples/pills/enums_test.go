package pills

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestEnums(t *testing.T) {
	t.Run("Pill", func(t *testing.T) {
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
					actual, ok := PillFromString(tC.upper)
					require.True(t, ok)
					require.Equal(t, tC.enum, actual)
					actual, ok = PillFromString(tC.lower)
					require.False(t, ok)
					require.Equal(t, Pill(0), actual)
				})
				t.Run(fmt.Sprintf("Case-insensitive lookup (idx: %d %s)", idx, tC.enum), func(t *testing.T) {
					enum, ok := PillFromStringIgnoreCase(tC.upper)
					require.True(t, ok)
					require.Equal(t, tC.enum, enum)
					enum, ok = PillFromStringIgnoreCase(tC.lower)
					require.True(t, ok)
					require.Equal(t, tC.enum, enum)
				})
			}
		})
		type testCase struct {
			enum            Pill
			from            string
			serialized      string
			wantStringerOut string
			invalid         bool
			shadowedBy      Pill
		}
		testCases := []testCase{
			{from: "", serialized: "", enum: Pill(-1), invalid: true, wantStringerOut: "Pill(-1)"},
			{from: "", serialized: "", enum: Pill(5), invalid: true, wantStringerOut: "Pill(4)"},
			{from: "PLACEBO", serialized: "PLACEBO", enum: Pill(0), wantStringerOut: "PLACEBO"},
			{from: "ASPIRIN", serialized: "ASPIRIN", enum: PillAspirin, wantStringerOut: "ASPIRIN"},
			{from: "IBUPROFEN", serialized: "IBUPROFEN", enum: PillIbuprofen, wantStringerOut: "IBUPROFEN"},
			{from: "PARACETAMOL", serialized: "PARACETAMOL", enum: PillParacetamol, wantStringerOut: "PARACETAMOL"},
			{from: "ACETAMINOPHEN", serialized: "ACETAMINOPHEN", enum: PillAcetaminophen, wantStringerOut: "ACETAMINOPHEN", shadowedBy: PillParacetamol},
			{from: "VITAMIN-C", serialized: "VITAMIN-C", enum: PillVitaminC, wantStringerOut: "VITAMIN-C"},
		}
		for idx, tC := range testCases {
			t.Run(fmt.Sprintf("Serializers (idx: %d %s)", idx, tC.enum), func(t *testing.T) {
				shadow := func(tC testCase) (after testCase) {
					if tC.shadowedBy == 0 {
						return tC
					}
					return testCase{
						enum:            tC.shadowedBy,
						from:            tC.from,
						serialized:      tC.shadowedBy.String(),
						wantStringerOut: tC.shadowedBy.String(),
						invalid:         tC.invalid,
					}
				}(tC)
				t.Run("binary", func(t *testing.T) {
					t.Run("MarhsalBinary", func(t *testing.T) {
						actual, err := tC.enum.MarshalBinary()
						if tC.invalid {
							require.Errorf(t, err, "got: %s", actual)
							return
						}
						require.NoError(t, err)
						require.Equal(t, shadow.serialized, string(actual))
					})
					t.Run("UnmarshalBinary", func(t *testing.T) {
						var enum Pill
						err := enum.UnmarshalBinary([]byte(tC.from))
						if tC.invalid {
							require.Error(t, err)
							return
						}
						require.NoError(t, err)
						require.Equal(t, shadow.enum, enum)
					})
				})
				t.Run("json", func(t *testing.T) {
					t.Run("MarhsalJSON", func(t *testing.T) {
						actual, err := tC.enum.MarshalJSON()
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
						var enum Pill
						jsonSerialized, err := json.Marshal(shadow.from)
						require.NoError(t, err)
						err = enum.UnmarshalJSON([]byte(jsonSerialized))
						if tC.invalid {
							require.Error(t, err)
							return
						}
						require.NoError(t, err)
						require.Equal(t, shadow.enum, enum)
					})
				})
				t.Run("text", func(t *testing.T) {
					t.Run("MarhsalText", func(t *testing.T) {
						actual, err := tC.enum.MarshalText()
						if tC.invalid {
							require.Error(t, err)
							return
						}
						require.NoError(t, err)
						require.Equal(t, shadow.serialized, string(actual))
					})
					t.Run("UnmarshalText", func(t *testing.T) {
						var enum Pill
						err := enum.UnmarshalText([]byte(tC.from))
						if tC.invalid {
							require.Error(t, err)
							return
						}
						require.NoError(t, err)
						require.Equal(t, tC.enum, enum)
					})
				})
				t.Run("yaml", func(t *testing.T) {
					t.Run("MarhsalYAML", func(t *testing.T) {
						j, err := tC.enum.MarshalYAML()
						if tC.invalid {
							require.Error(t, err)
							return
						}
						require.NoError(t, err)
						require.Equal(t, shadow.serialized, j)
					})
					t.Run("UnmarshalYAML", func(t *testing.T) {
						var enum Pill
						err := enum.UnmarshalYAML(&yaml.Node{Kind: yaml.ScalarNode, Value: tC.serialized})
						if tC.invalid {
							require.Error(t, err)
							return
						}
						require.NoError(t, err)
						require.Equal(t, shadow.enum, enum)
					})
				})
				t.Run("sql", func(t *testing.T) {
					t.Run("Value", func(t *testing.T) {
						j, err := tC.enum.Value()
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
							var enum Pill
							err := enum.Scan(v)
							if tC.invalid {
								require.Error(t, err)
								return
							}
							require.NoError(t, err)
							require.Equal(t, shadow.enum, enum)
						}
					})
					t.Run("Scan <nil>", func(t *testing.T) {
						var value interface{} = nil
						var enum Pill
						err := enum.Scan(value)
						require.Error(t, err)
						require.Zero(t, enum)
					})
				})
			})
		}
	})
	t.Run("PillUndefined", func(t *testing.T) {
		// TODO: this
	})
	t.Run("PillAutoStripPrefix", func(t *testing.T) {
		// TODO:
		// Assert that
	})
}

type stringer struct{ v string }

func (s stringer) String() string { return s.v }
