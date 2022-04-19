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
			t.Run("return copies", func(t *testing.T) {
				utils.AssertNotSamePointer(t, _BookingStateStrings, BookingStateStrings())
				utils.AssertNotSamePointer(t, _BookingStateValues, BookingStateValues())
			})
			t.Run("additional data", func(t *testing.T) {
				descriptions := []string{"The booking was created successfully", "The booking was not available", "The booking failed", "The booking was canceled", "The booking was not found", "The booking was deleted"}
				for _, enum := range BookingStateValues() {
					require.Equal(t, descriptions[enum], utils.Must(enum.GetDescription()))
				}
				t.Run("for invalid enum", func(t *testing.T) {
					desc, ok := BookingState(999).GetDescription()
					require.False(t, ok)
					require.Equal(t, "BookingState(999).Description", desc)
				})
			})
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
					enum, ok := BookingStateFromString(tC.upper)
					require.True(t, ok)
					require.Equal(t, tC.enum, enum)
					require.Equal(t, tC.canonical, utils.Must(enum.GetDescription()))
					enum, ok = BookingStateFromString(tC.lower)
					if tC.lower == tC.upper {
						require.True(t, ok)
						require.Equal(t, tC.enum, enum)
						require.Equal(t, tC.canonical, utils.Must(enum.GetDescription()))
					} else {
						require.False(t, ok)
						require.Equal(t, BookingState(0), enum)
					}
				})
				t.Run(fmt.Sprintf("Case-insensitive lookup (idx: %d %s)", idx, tC.enum), func(t *testing.T) {
					enum, ok := BookingStateFromStringIgnoreCase(tC.upper)
					require.True(t, ok)
					require.Equal(t, tC.enum, enum)
					require.Equal(t, tC.canonical, utils.Must(enum.GetDescription()))
					enum, ok = BookingStateFromStringIgnoreCase(tC.lower)
					require.True(t, ok)
					require.Equal(t, tC.enum, enum)
					require.Equal(t, tC.canonical, utils.Must(enum.GetDescription()))
				})
			}
		})
		t.Run("Missing Serializers", func(t *testing.T) {
			utils.AssertMissingSerializationInterfacesFor[BookingState](t, []string{"binary", "gql", "json", "sql", "text"})
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
				serializers := []string{"yaml"}
				utils.AssertSerializationInterfacesFor[BookingState](t, idx, tC, cfg, serializers)
			}
		})
	})
	t.Run("BookingStateWithConfig", func(t *testing.T) {
		// TODO: Extract: Overwriting configuration on a type spec level deserves its own use case package
		t.Run("Value Sets", func(t *testing.T) {
			require.Equal(t,
				[]string{"Created", "Unavailable", "Failed", "Canceled", "NotFound", "Deleted"},
				BookingStateWithConfigStrings())
			require.Equal(t,
				[]BookingStateWithConfig{0, 1, 2, 3, 4, 5},
				BookingStateWithConfigValues())
			t.Run("return copies", func(t *testing.T) {
				utils.AssertNotSamePointer(t, _BookingStateWithConfigStrings, BookingStateWithConfigStrings())
				utils.AssertNotSamePointer(t, _BookingStateWithConfigValues, BookingStateWithConfigValues())
			})
			t.Run("Misses Ent Interface", func(t *testing.T) {
				_, ok := ((interface{})(BookingStateWithConfig(0))).(interface{ Values() []string })
				require.False(t, ok)
			})
		})
		t.Run("Lookup", func(t *testing.T) {
			type testCase struct {
				enum      BookingStateWithConfig
				upper     string
				lower     string
				canonical string
			}
			testCases := []testCase{
				{BookingStateWithConfig(0), "Created", "created", "The booking was created successfully"},
				{BookingStateWithConfig(1), "Unavailable", "unavailable", "The booking was not available"},
				{BookingStateWithConfig(2), "Failed", "failed", "The booking failed"},
				{BookingStateWithConfig(3), "Canceled", "canceled", "The booking was canceled"},
				{BookingStateWithConfig(4), "NotFound", "notfound", "The booking was not found"},
				{BookingStateWithConfig(5), "Deleted", "deleted", "The booking was deleted"},
			}
			for idx, tC := range testCases {
				t.Run(fmt.Sprintf("Case-sensitive lookup (idx: %d %s)", idx, tC.enum), func(t *testing.T) {
					enum, ok := BookingStateWithConfigFromString(tC.upper)
					require.True(t, ok)
					require.Equal(t, tC.enum, enum)
					require.Equal(t, tC.canonical, utils.Must(enum.GetDescription()))
					enum, ok = BookingStateWithConfigFromString(tC.lower)
					if tC.lower == tC.upper {
						require.True(t, ok)
						require.Equal(t, tC.enum, enum)
						require.Equal(t, tC.canonical, utils.Must(enum.GetDescription()))
					} else {
						require.False(t, ok)
						require.Equal(t, BookingStateWithConfig(0), enum)
					}
				})
				t.Run(fmt.Sprintf("Case-insensitive lookup (idx: %d %s)", idx, tC.enum), func(t *testing.T) {
					enum, ok := BookingStateWithConfigFromStringIgnoreCase(tC.upper)
					require.True(t, ok)
					require.Equal(t, tC.enum, enum)
					require.Equal(t, tC.canonical, utils.Must(enum.GetDescription()))
					enum, ok = BookingStateWithConfigFromStringIgnoreCase(tC.lower)
					require.True(t, ok)
					require.Equal(t, tC.enum, enum)
					require.Equal(t, tC.canonical, utils.Must(enum.GetDescription()))
				})
			}
		})
		t.Run("Missing Serializers", func(t *testing.T) {
			utils.AssertMissingSerializationInterfacesFor[BookingStateWithConfig](t, []string{"binary", "gql", "sql", "text"})
		})
		t.Run("Serialization", func(t *testing.T) {
			cfg := utils.TestConfig{}
			toPtr := utils.ToPointer[BookingStateWithConfig]
			testCases := []utils.TestCase{
				{From: "UNKNOWN", Enum: toPtr(6), Expected: utils.Expected{AsSerialized: "BookingStateWithConfig(6)", IsInvalid: true}},
				{From: "Created", Enum: toPtr(0), Expected: utils.Expected{AsSerialized: "Created"}},
				{From: "Unavailable", Enum: toPtr(1), Expected: utils.Expected{AsSerialized: "Unavailable"}},
				{From: "Failed", Enum: toPtr(2), Expected: utils.Expected{AsSerialized: "Failed"}},
				{From: "Canceled", Enum: toPtr(3), Expected: utils.Expected{AsSerialized: "Canceled"}},
				{From: "NotFound", Enum: toPtr(4), Expected: utils.Expected{AsSerialized: "NotFound"}},
				{From: "Deleted", Enum: toPtr(5), Expected: utils.Expected{AsSerialized: "Deleted"}},
			}
			for idx, tC := range testCases {
				serializers := []string{"json", "yaml"}
				utils.AssertSerializationInterfacesFor[BookingStateWithConfig](t, idx, tC, cfg, serializers)
			}
		})
	})
	t.Run("BookingStateWithConstants", func(t *testing.T) {
		t.Run("Value Sets", func(t *testing.T) {
			require.Equal(t,
				[]string{"Created", "Unavailable", "Failed", "Canceled", "NotFound", "Deleted"},
				BookingStateWithConstantsStrings())
			require.Equal(t,
				[]BookingStateWithConstants{0, 1, 2, 3, 4, 5},
				BookingStateWithConstantsValues())
			t.Run("return copies", func(t *testing.T) {
				utils.AssertNotSamePointer(t, _BookingStateWithConstantsStrings, BookingStateWithConstantsStrings())
				utils.AssertNotSamePointer(t, _BookingStateWithConstantsValues, BookingStateWithConstantsValues())
			})
			t.Run("Ent Interface", func(t *testing.T) {
				require.Equal(t,
					[]string{"Created", "Unavailable", "Failed", "Canceled", "NotFound", "Deleted"},
					BookingStateWithConstants(0).Values())
			})
		})
		t.Run("Lookup", func(t *testing.T) {
			type testCase struct {
				enum      BookingStateWithConstants
				upper     string
				lower     string
				canonical string
			}
			testCases := []testCase{
				{BookingStateWithConstants(0), "Created", "created", "The booking was created successfully"},
				{BookingStateWithConstants(1), "Unavailable", "unavailable", "The booking was not available"},
				{BookingStateWithConstants(2), "Failed", "failed", "The booking failed"},
				{BookingStateWithConstants(3), "Canceled", "canceled", "The booking was canceled"},
				{BookingStateWithConstants(4), "NotFound", "notfound", "The booking was not found"},
				{BookingStateWithConstants(5), "Deleted", "deleted", "The booking was deleted"},
			}
			for idx, tC := range testCases {
				t.Run(fmt.Sprintf("Case-sensitive lookup (idx: %d %s)", idx, tC.enum), func(t *testing.T) {
					enum, ok := BookingStateWithConstantsFromString(tC.upper)
					require.True(t, ok)
					require.Equal(t, tC.enum, enum)
					require.Equal(t, tC.canonical, utils.Must(enum.GetDescription()))
					enum, ok = BookingStateWithConstantsFromString(tC.lower)
					if tC.lower == tC.upper {
						require.True(t, ok)
						require.Equal(t, tC.enum, enum)
						require.Equal(t, tC.canonical, utils.Must(enum.GetDescription()))
					} else {
						require.False(t, ok)
						require.Equal(t, BookingStateWithConstants(0), enum)
					}
				})
				t.Run(fmt.Sprintf("Case-insensitive lookup (idx: %d %s)", idx, tC.enum), func(t *testing.T) {
					enum, ok := BookingStateWithConstantsFromStringIgnoreCase(tC.upper)
					require.True(t, ok)
					require.Equal(t, tC.enum, enum)
					require.Equal(t, tC.canonical, utils.Must(enum.GetDescription()))
					enum, ok = BookingStateWithConstantsFromStringIgnoreCase(tC.lower)
					require.True(t, ok)
					require.Equal(t, tC.enum, enum)
					require.Equal(t, tC.canonical, utils.Must(enum.GetDescription()))
				})
			}
		})
		t.Run("Missing Serializers", func(t *testing.T) {
			utils.AssertMissingSerializationInterfacesFor[BookingStateWithConstants](t, []string{"binary", "gql", "json", "sql", "text"})
		})
		t.Run("Serialization", func(t *testing.T) {
			cfg := utils.TestConfig{}
			toPtr := utils.ToPointer[BookingStateWithConstants]
			testCases := []utils.TestCase{
				{From: "UNKNOWN", Enum: toPtr(6), Expected: utils.Expected{AsSerialized: "BookingStateWithConstants(6)", IsInvalid: true}},
				{From: "Created", Enum: toPtr(0), Expected: utils.Expected{AsSerialized: "Created"}},
				{From: "Unavailable", Enum: toPtr(1), Expected: utils.Expected{AsSerialized: "Unavailable"}},
				{From: "Failed", Enum: toPtr(2), Expected: utils.Expected{AsSerialized: "Failed"}},
				{From: "Canceled", Enum: toPtr(3), Expected: utils.Expected{AsSerialized: "Canceled"}},
				{From: "NotFound", Enum: toPtr(4), Expected: utils.Expected{AsSerialized: "NotFound"}},
				{From: "Deleted", Enum: toPtr(5), Expected: utils.Expected{AsSerialized: "Deleted"}},
			}
			for idx, tC := range testCases {
				serializers := []string{"yaml"}
				utils.AssertSerializationInterfacesFor[BookingStateWithConstants](t, idx, tC, cfg, serializers)
			}
		})
	})
}
