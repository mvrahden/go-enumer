package project

import (
	"testing"

	"github.com/mvrahden/go-enumer/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestEnums(t *testing.T) {
	t.Run("AccountState", func(t *testing.T) {
		t.Run("Value Sets", func(t *testing.T) {
			t.Run("return copies", func(t *testing.T) {
				utils.AssertNotSamePointer(t, _AccountStateStrings, AccountStateStrings())
				utils.AssertNotSamePointer(t, _AccountStateValues, AccountStateValues())
			})
		})
		t.Run("Missing Serializers", func(t *testing.T) {
			utils.AssertMissingSerializationInterfacesFor[AccountState](t, []string{"binary", "gql", "text", "yaml", "yaml.v3"})
		})
		t.Run("Serialization", func(t *testing.T) {
			cfg := utils.TestConfig{}
			toPtr := utils.ToPointer[AccountState]
			testCases := []utils.TestCase{
				// hint: this 1st case is invalid upon deserialization,
				// but valid upon serialization (as it is the default value
				// but does not support "undefined")
				{From: "", Enum: toPtr(0), Expected: utils.Expected{AsSerialized: "STAGED", IsInvalid: true, IsDefault: true}},
				{From: "NOT_A_ACCOUNT_STATE", Enum: toPtr(5), Expected: utils.Expected{AsSerialized: "AccountState(5)", IsInvalid: true}},
				{From: "STAGED", Enum: toPtr(0), Expected: utils.Expected{AsSerialized: "STAGED"}},
				{From: "PROVISIONED", Enum: toPtr(1), Expected: utils.Expected{AsSerialized: "PROVISIONED"}},
				{From: "ACTIVATED", Enum: toPtr(2), Expected: utils.Expected{AsSerialized: "ACTIVATED"}},
				{From: "DEACTIVATED", Enum: toPtr(3), Expected: utils.Expected{AsSerialized: "DEACTIVATED"}},
				{From: "DEPROVISIONED", Enum: toPtr(4), Expected: utils.Expected{AsSerialized: "DEPROVISIONED"}},
			}
			for idx, tC := range testCases {
				serializers := []string{"json", "sql"}
				utils.AssertSerializationInterfacesFor[AccountState](t, idx, tC, cfg, serializers)
			}
		})
	})
	t.Run("CountryCode", func(t *testing.T) {
		t.Run("Value Sets", func(t *testing.T) {
			t.Run("return copies", func(t *testing.T) {
				utils.AssertNotSamePointer(t, _CountryCodeStrings, CountryCodeStrings())
				utils.AssertNotSamePointer(t, _CountryCodeValues, CountryCodeValues())
			})
		})
		t.Run("Additional Data", func(t *testing.T) {
			require.Equal(t, "United States", CountryCode(229).GetCountryName())
			require.Equal(t, "US", CountryCode(229).GetIso2LetterCode())
			require.Equal(t, "1", CountryCode(229).GetCountryCode())
			require.Equal(t, uint32(310232863), CountryCode(229).GetPopulation())
			require.Equal(t, uint32(9629091), CountryCode(229).GetAreaInSquareKilometer())
			require.Equal(t, float64(1672), CountryCode(229).GetGdpInBillion())
			t.Run("panics for invalid enum", func(t *testing.T) {
				// hint: during runtime it is necessary for us to have valid enums.
				// Panics can help to detect unhandled invalid enums early in your development process.
				// Unhandled invalid enums can sneak their way in easily e.g. due to improper
				// intialization of an enum variable or missing validation.
				require.PanicsWithError(t,
					"Forbidden access to additional enum data of \"CountryCode(0)\". err: not a valid enum",
					func() {
						CountryCode(0).GetCountryName()
					})
			})
		})
		t.Run("Serialization", func(t *testing.T) {
			cfg := utils.TestConfig{}
			toPtr := utils.ToPointer[CountryCode]
			testCases := []utils.TestCase{
				{From: "", Enum: toPtr(0), Expected: utils.Expected{AsSerialized: "CountryCode(0)", IsInvalid: true}},
				{From: "NOT_A_COUNTRY_CODE", Enum: toPtr(241), Expected: utils.Expected{AsSerialized: "CountryCode(241)", IsInvalid: true}},
				{From: "AFG", Enum: toPtr(1), Expected: utils.Expected{AsSerialized: "AFG"}},
				{From: "ALB", Enum: toPtr(2), Expected: utils.Expected{AsSerialized: "ALB"}},
				{From: "ZWE", Enum: toPtr(240), Expected: utils.Expected{AsSerialized: "ZWE"}},
				{From: "DEU", Enum: toPtr(78), Expected: utils.Expected{AsSerialized: "DEU"}},
				{From: "USA", Enum: toPtr(229), Expected: utils.Expected{AsSerialized: "USA"}},
			}
			for idx, tC := range testCases {
				serializers := []string{"binary", "gql", "json", "sql", "text", "yaml"}
				utils.AssertSerializationInterfacesFor[CountryCode](t, idx, tC, cfg, serializers)
			}
		})
	})
	t.Run("Currency", func(t *testing.T) {
		t.Run("Value Sets", func(t *testing.T) {
			t.Run("return copies", func(t *testing.T) {
				utils.AssertNotSamePointer(t, _CurrencyStrings, CurrencyStrings())
				utils.AssertNotSamePointer(t, _CurrencyValues, CurrencyValues())
			})
		})
		t.Run("Serialization", func(t *testing.T) {
			cfg := utils.TestConfig{}
			toPtr := utils.ToPointer[Currency]
			testCases := []utils.TestCase{
				{From: "", Enum: toPtr(0), Expected: utils.Expected{AsSerialized: "Currency(0)", IsInvalid: true}},
				{From: "NOT_A_CURRENCY", Enum: toPtr(6), Expected: utils.Expected{AsSerialized: "Currency(6)", IsInvalid: true}},
				{From: "USD", Enum: toPtr(1), Expected: utils.Expected{AsSerialized: "USD"}},
				{From: "EUR", Enum: toPtr(2), Expected: utils.Expected{AsSerialized: "EUR"}},
				{From: "JPY", Enum: toPtr(3), Expected: utils.Expected{AsSerialized: "JPY"}},
				{From: "GBP", Enum: toPtr(4), Expected: utils.Expected{AsSerialized: "GBP"}},
				{From: "AUD", Enum: toPtr(5), Expected: utils.Expected{AsSerialized: "AUD"}},
			}
			for idx, tC := range testCases {
				serializers := []string{"binary", "gql", "json", "sql", "text", "yaml"}
				utils.AssertSerializationInterfacesFor[Currency](t, idx, tC, cfg, serializers)
			}
		})
	})
	t.Run("Timezone", func(t *testing.T) {
		t.Run("Value Sets", func(t *testing.T) {
			t.Run("return copies", func(t *testing.T) {
				utils.AssertNotSamePointer(t, _TimezoneStrings, TimezoneStrings())
				utils.AssertNotSamePointer(t, _TimezoneValues, TimezoneValues())
			})
		})
		t.Run("Serialization", func(t *testing.T) {
			cfg := utils.TestConfig{}
			toPtr := utils.ToPointer[Timezone]
			testCases := []utils.TestCase{
				{From: "", Enum: toPtr(0), Expected: utils.Expected{AsSerialized: "Timezone(0)", IsInvalid: true}},
				{From: "NOT_A_TIMEZONE", Enum: toPtr(425), Expected: utils.Expected{AsSerialized: "Timezone(425)", IsInvalid: true}},
				{From: "Asia/Kabul", Enum: toPtr(1), Expected: utils.Expected{AsSerialized: "Asia/Kabul"}},
				{From: "Asia/Beirut", Enum: toPtr(217), Expected: utils.Expected{AsSerialized: "Asia/Beirut"}},
				{From: "Europe/London", Enum: toPtr(TimezoneEuropeLondon), Expected: utils.Expected{AsSerialized: "Europe/London"}},
				{From: "America/New_York", Enum: toPtr(400), Expected: utils.Expected{AsSerialized: "America/New_York"}},
				{From: "Europe/Mariehamn", Enum: toPtr(424), Expected: utils.Expected{AsSerialized: "Europe/Mariehamn"}},
			}
			for idx, tC := range testCases {
				serializers := []string{"binary", "gql", "json", "sql", "text", "yaml"}
				utils.AssertSerializationInterfacesFor[Timezone](t, idx, tC, cfg, serializers)
			}
		})
	})
	t.Run("UserRole", func(t *testing.T) {
		t.Run("Value Sets", func(t *testing.T) {
			t.Run("return copies", func(t *testing.T) {
				utils.AssertNotSamePointer(t, _UserRoleStrings, UserRoleStrings())
				utils.AssertNotSamePointer(t, _UserRoleValues, UserRoleValues())
			})
		})
		t.Run("Serialization", func(t *testing.T) {
			cfg := utils.TestConfig{SupportUndefined: true}
			toPtr := utils.ToPointer[UserRole]
			testCases := []utils.TestCase{
				{From: "", Enum: toPtr(0), Expected: utils.Expected{AsSerialized: "standard"}},
				{From: "NOT_A_USER_ROLE", Enum: toPtr(4), Expected: utils.Expected{AsSerialized: "UserRole(4)", IsInvalid: true}},
				{From: "standard", Enum: toPtr(UserRoleStandard), Expected: utils.Expected{AsSerialized: "standard"}},
				{From: "editor", Enum: toPtr(UserRoleEditor), Expected: utils.Expected{AsSerialized: "editor"}},
				{From: "reviewer", Enum: toPtr(UserRoleReviewer), Expected: utils.Expected{AsSerialized: "reviewer"}},
				{From: "admin", Enum: toPtr(UserRoleAdmin), Expected: utils.Expected{AsSerialized: "admin"}},
			}
			for idx, tC := range testCases {
				serializers := []string{"binary", "gql", "json", "sql", "text", "yaml"}
				utils.AssertSerializationInterfacesFor[UserRole](t, idx, tC, cfg, serializers)
			}
		})
	})
}
