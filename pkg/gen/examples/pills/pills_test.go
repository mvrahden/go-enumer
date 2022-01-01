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
			[]string{"Placebo", "Aspirin", "Ibuprofen", "Acetaminophen"},
			PillStrings())
		require.Equal(t,
			[]Pill{PillPlacebo, PillAspirin, PillIbuprofen, PillAcetaminophen},
			PillValues())
		t.Run("", func(t *testing.T) {

		})
	})
	t.Run("Serialization", func(t *testing.T) {
		testCases := []struct {
			p          Pill
			serialized string
			invalid    bool
		}{
			{serialized: "", p: Pill(-1), invalid: true},
			{serialized: "Placebo"},
			{serialized: "Aspirin", p: PillAspirin},
			{serialized: "Ibuprofen", p: PillIbuprofen},
			{serialized: "Acetaminophen", p: PillAcetaminophen},
		}
		for _, tC := range testCases {
			t.Run("binary", func(t *testing.T) {
				t.Run("MarhsalBinary", func(t *testing.T) {
					j, err := tC.p.MarshalBinary()
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, tC.serialized, string(j))
				})
				t.Run("UnmarshalBinary", func(t *testing.T) {
					p := tC.p
					err := p.UnmarshalBinary([]byte(tC.serialized))
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, tC.p, p)
				})
			})
			t.Run("json", func(t *testing.T) {
				jsonSerialized := fmt.Sprintf("%q", tC.serialized)
				t.Run("MarhsalJSON", func(t *testing.T) {
					j, err := tC.p.MarshalJSON()
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, jsonSerialized, string(j))
				})
				t.Run("UnmarshalJSON", func(t *testing.T) {
					p := tC.p
					err := p.UnmarshalJSON([]byte(jsonSerialized))
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, tC.p, p)
				})
			})
			t.Run("text", func(t *testing.T) {
				t.Run("MarhsalText", func(t *testing.T) {
					j, err := tC.p.MarshalText()
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, tC.serialized, string(j))
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
					require.Equal(t, tC.serialized, j)
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
					require.Equal(t, tC.p, p)
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
					require.Equal(t, tC.serialized, j)
				})
				t.Run("Scan", func(t *testing.T) {
					p := tC.p
					err := p.Scan(tC.serialized)
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, tC.p, p)
				})
			})
		}
	})
}