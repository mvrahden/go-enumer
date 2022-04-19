package project

import (
	"testing"

	"github.com/mvrahden/go-enumer/pkg/utils"
)

func TestEnums(t *testing.T) {
	t.Run("UserRole", func(t *testing.T) {
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
	t.Run("Timezone", func(t *testing.T) {
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
	t.Run("Currency", func(t *testing.T) {
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
	t.Run("AccountState", func(t *testing.T) {
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
}
