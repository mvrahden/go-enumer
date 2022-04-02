package booking

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGreetings(t *testing.T) {
	t.Run("Value Sets", func(t *testing.T) {
		require.Equal(t,
			[]string{"Created", "Unavailable", "Failed", "Canceled", "NotFound", "Deleted"},
			BookingStateStrings())
		require.Equal(t,
			[]BookingState{0, 1, 2, 3, 4, 5},
			BookingStateValues())
		t.Run("Ent Interface", func(t *testing.T) {
			require.Equal(t,
				[]string{"Created", "Unavailable", "Failed", "Canceled", "NotFound", "Deleted"},
				BookingState(0).Values())
		})
	})
	t.Run("Lookup", func(t *testing.T) {
		type testCase struct {
			enum  BookingState
			upper string
			lower string
		}
		testCases := []testCase{
			{BookingState(0), "Created", "created"},
			{BookingState(1), "Unavailable", "unavailable"},
			{BookingState(2), "Failed", "failed"},
			{BookingState(3), "Canceled", "canceled"},
			{BookingState(4), "NotFound", "notfound"},
			{BookingState(5), "Deleted", "deleted"},
		}
		for idx, tC := range testCases {
			t.Run(fmt.Sprintf("Case-sensitive lookup (idx: %d %s)", idx, tC.enum), func(t *testing.T) {
				p, ok := BookingStateFromString(tC.upper)
				require.True(t, ok)
				require.Equal(t, tC.enum, p)
				p, ok = BookingStateFromString(tC.lower)
				if tC.lower == tC.upper {
					require.True(t, ok)
					require.Equal(t, tC.enum, p)
				} else {
					require.False(t, ok)
					require.Equal(t, BookingState(0), p)
				}
			})
			t.Run(fmt.Sprintf("Case-insensitive lookup (idx: %d %s)", idx, tC.enum), func(t *testing.T) {
				p, ok := BookingStateFromStringIgnoreCase(tC.upper)
				require.True(t, ok)
				require.Equal(t, tC.enum, p)
				p, ok = BookingStateFromStringIgnoreCase(tC.lower)
				require.True(t, ok)
				require.Equal(t, tC.enum, p)
			})
		}
	})
	testCases := []struct {
		from       string
		enum       BookingState
		serialized string
		invalid    bool
		stringer   string
		canonical  string
	}{
		{from: "", serialized: "Created", enum: BookingState(0), invalid: false, stringer: "Created", canonical: "The booking was created successfully"},
		{from: "BookingState(6)", serialized: "BookingState(6)", enum: BookingState(6), invalid: true, stringer: "BookingState(6)", canonical: "BookingState(6)"},
		{from: "Created", serialized: "Created", enum: 0, stringer: "Created", canonical: "The booking was created successfully"},
		{from: "Unavailable", serialized: "Unavailable", enum: 1, stringer: "Unavailable", canonical: "The booking was not available"},
		{from: "Failed", serialized: "Failed", enum: 2, stringer: "Failed", canonical: "The booking failed"},
		{from: "Canceled", serialized: "Canceled", enum: 3, stringer: "Canceled", canonical: "The booking was canceled"},
		{from: "NotFound", serialized: "NotFound", enum: 4, stringer: "NotFound", canonical: "The booking was not found"},
		{from: "Deleted", serialized: "Deleted", enum: 5, stringer: "Deleted", canonical: "The booking was deleted"},
	}
	for idx, tC := range testCases {
		t.Run(fmt.Sprintf("Serializers (idx: %d %s)", idx, tC.enum), func(t *testing.T) {
			t.Run("Stringer", func(t *testing.T) {
				require.Equal(t, tC.stringer, tC.enum.String())
			})
			t.Run("CanonicalValue", func(t *testing.T) {
				require.Equal(t, tC.canonical, tC.enum.CanonicalValue())
			})
			t.Run("yaml", func(t *testing.T) {
				t.Run("MarhsalYAML", func(t *testing.T) {
					j, err := tC.enum.MarshalYAML()
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, tC.serialized, j)
				})
				t.Run("UnmarshalYAML", func(t *testing.T) {
					var g BookingState
					err := g.UnmarshalYAML(func(i interface{}) error {
						return json.Unmarshal([]byte("\""+tC.serialized+"\""), i)
					})
					if tC.invalid {
						require.Error(t, err)
						return
					}
					require.NoError(t, err)
					require.Equal(t, tC.enum, g)
				})
			})
		})
	}
}
