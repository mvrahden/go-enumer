package booking

import (
	"fmt"
	"testing"

	"github.com/mvrahden/go-enumer/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestEnums(t *testing.T) {
	t.Run("BookingState", func(t *testing.T) {
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
				enum      BookingState
				upper     string
				lower     string
				canonical string
			}
			testCases := []testCase{
				{BookingState(0), "Created", "created", "The booking was created successfully"},
				{BookingState(1), "Unavailable", "unavailable", "The booking was not available"},
				{BookingState(2), "Failed", "failed", "The booking failed"},
				{BookingState(3), "Canceled", "canceled", "The booking was canceled"},
				{BookingState(4), "NotFound", "notfound", "The booking was not found"},
				{BookingState(5), "Deleted", "deleted", "The booking was deleted"},
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
						require.Equal(t, tC.canonical, p.CanonicalValue())
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
					require.Equal(t, tC.canonical, p.CanonicalValue())
				})
			}
		})
		t.Run("Serialization", func(t *testing.T) {
			cfg := utils.TestConfig{}
			toPtr := utils.ToPointer[BookingState]
			testCases := []utils.TestCase{
				{From: "UNKNOWN", Enum: toPtr(6), Expected: utils.Expected{AsSerialized: "BookingState(6)", IsInvalid: true}},
				{From: "Created", Enum: toPtr(0), Expected: utils.Expected{AsSerialized: "Created"}},
				{From: "Unavailable", Enum: toPtr(1), Expected: utils.Expected{AsSerialized: "Unavailable"}},
				{From: "Failed", Enum: toPtr(2), Expected: utils.Expected{AsSerialized: "Failed"}},
				{From: "Canceled", Enum: toPtr(3), Expected: utils.Expected{AsSerialized: "Canceled"}},
				{From: "NotFound", Enum: toPtr(4), Expected: utils.Expected{AsSerialized: "NotFound"}},
				{From: "Deleted", Enum: toPtr(5), Expected: utils.Expected{AsSerialized: "Deleted"}},
			}
			for idx, tC := range testCases {
				t.Run(fmt.Sprintf("Serializers (idx: %d from %q)", idx, tC.From), func(t *testing.T) {
					utils.AssertSerializers[BookingState](t, tC, "yaml")
				})
				t.Run(fmt.Sprintf("Deserializers (idx: %d from %q)", idx, tC.From), func(t *testing.T) {
					utils.AssertDeserializers(t, tC, cfg, "yaml", utils.ZeroValuer[BookingState])
				})
			}
		})
	})
	t.Run("BookingStateWithConstants", func(t *testing.T) {
		// TODO:
		// - assert that type spec is derived from CSV and NOT from constants
	})
	t.Run("BookingStateWithConfig", func(t *testing.T) {
		// TODO:
		// This use case needs extraction...
	})
}
