package planets

import (
	"testing"

	"github.com/mvrahden/go-enumer/pkg/utils"
)

func TestEnums(t *testing.T) {
	t.Run("Planet", func(t *testing.T) {
		t.Run("Serialization", func(t *testing.T) {
			cfg := utils.TestConfig{}
			toPtr := utils.ToPointer[Planet]
			testCases := []utils.TestCase{
				{From: "", Enum: toPtr(0), Expected: utils.Expected{AsSerialized: "Planet(0)", IsInvalid: true}},
				{From: "UNKNOWN", Enum: toPtr(9), Expected: utils.Expected{AsSerialized: "Planet(9)", IsInvalid: true}},
				{From: "Mars", Enum: toPtr(PlanetMars), Expected: utils.Expected{AsSerialized: "Mars"}},
				{From: "Pluto", Enum: toPtr(PlanetPluto), Expected: utils.Expected{AsSerialized: "Pluto"}},
				{From: "Venus", Enum: toPtr(PlanetVenus), Expected: utils.Expected{AsSerialized: "Venus"}},
				{From: "Mercury", Enum: toPtr(PlanetMercury), Expected: utils.Expected{AsSerialized: "Mercury"}},
				{From: "Jupiter", Enum: toPtr(PlanetJupiter), Expected: utils.Expected{AsSerialized: "Jupiter"}},
				{From: "Saturn", Enum: toPtr(PlanetSaturn), Expected: utils.Expected{AsSerialized: "Saturn"}},
				{From: "Uranus", Enum: toPtr(PlanetUranus), Expected: utils.Expected{AsSerialized: "Uranus"}},
				{From: "Neptune", Enum: toPtr(PlanetNeptune), Expected: utils.Expected{AsSerialized: "Neptune"}},
			}
			for idx, tC := range testCases {
				serializers := []string{"binary", "gql", "json", "sql", "text", "yaml"}
				utils.AssertSerializationInterfacesFor[Planet](t, idx, tC, cfg, serializers)
			}
		})
	})
	t.Run("PlanetWithDefault", func(t *testing.T) {
		t.Run("Serialization", func(t *testing.T) {
			cfg := utils.TestConfig{HasDefault: true}
			toPtr := utils.ToPointer[PlanetWithDefault]
			testCases := []utils.TestCase{
				// hint: this 1st case is invalid upon deserialization,
				// but valid upon serialization (as it is the default value
				// but does not support "undefined")
				{From: "", Enum: toPtr(0), Expected: utils.Expected{AsSerialized: "Earth", IsInvalid: true}},
				{From: "PlanetWithDefault(9)", Enum: toPtr(9), Expected: utils.Expected{AsSerialized: "PlanetWithDefault(9)", IsInvalid: true}},
				{From: "Earth", Enum: toPtr(0), Expected: utils.Expected{AsSerialized: "Earth"}},
				{From: "Mars", Enum: toPtr(PlanetWithDefaultMars), Expected: utils.Expected{AsSerialized: "Mars"}},
				{From: "Pluto", Enum: toPtr(PlanetWithDefaultPluto), Expected: utils.Expected{AsSerialized: "Pluto"}},
				{From: "Venus", Enum: toPtr(PlanetWithDefaultVenus), Expected: utils.Expected{AsSerialized: "Venus"}},
				{From: "Mercury", Enum: toPtr(PlanetWithDefaultMercury), Expected: utils.Expected{AsSerialized: "Mercury"}},
				{From: "Jupiter", Enum: toPtr(PlanetWithDefaultJupiter), Expected: utils.Expected{AsSerialized: "Jupiter"}},
				{From: "Saturn", Enum: toPtr(PlanetWithDefaultSaturn), Expected: utils.Expected{AsSerialized: "Saturn"}},
				{From: "Uranus", Enum: toPtr(PlanetWithDefaultUranus), Expected: utils.Expected{AsSerialized: "Uranus"}},
				{From: "Neptune", Enum: toPtr(PlanetWithDefaultNeptune), Expected: utils.Expected{AsSerialized: "Neptune"}},
			}
			for idx, tC := range testCases {
				serializers := []string{"binary", "gql", "json", "sql", "text", "yaml"}
				utils.AssertSerializationInterfacesFor[PlanetWithDefault](t, idx, tC, cfg, serializers)
			}
		})
	})
	t.Run("PlanetSupportUndefined", func(t *testing.T) {
		t.Run("Serialization", func(t *testing.T) {
			toPtr := utils.ToPointer[PlanetSupportUndefined]
			cfg := utils.TestConfig{SupportUndefined: true}
			testCases := []utils.TestCase{
				{From: "", Enum: toPtr(0), Expected: utils.Expected{AsSerialized: "", IsInvalid: false}},
				{From: "PlanetSupportUndefined(9)", Enum: toPtr(9), Expected: utils.Expected{AsSerialized: "PlanetSupportUndefined(9)", IsInvalid: true}},
				{From: "Mars", Enum: toPtr(PlanetSupportUndefinedMars), Expected: utils.Expected{AsSerialized: "Mars"}},
				{From: "Pluto", Enum: toPtr(PlanetSupportUndefinedPluto), Expected: utils.Expected{AsSerialized: "Pluto"}},
				{From: "Venus", Enum: toPtr(PlanetSupportUndefinedVenus), Expected: utils.Expected{AsSerialized: "Venus"}},
				{From: "Mercury", Enum: toPtr(PlanetSupportUndefinedMercury), Expected: utils.Expected{AsSerialized: "Mercury"}},
				{From: "Jupiter", Enum: toPtr(PlanetSupportUndefinedJupiter), Expected: utils.Expected{AsSerialized: "Jupiter"}},
				{From: "Saturn", Enum: toPtr(PlanetSupportUndefinedSaturn), Expected: utils.Expected{AsSerialized: "Saturn"}},
				{From: "Uranus", Enum: toPtr(PlanetSupportUndefinedUranus), Expected: utils.Expected{AsSerialized: "Uranus"}},
				{From: "Neptune", Enum: toPtr(PlanetSupportUndefinedNeptune), Expected: utils.Expected{AsSerialized: "Neptune"}},
			}
			for idx, tC := range testCases {
				serializers := []string{"binary", "gql", "json", "sql", "text", "yaml"}
				utils.AssertSerializationInterfacesFor[PlanetSupportUndefined](t, idx, tC, cfg, serializers)
			}
		})
	})
	t.Run("PlanetSupportUndefinedWithDefault", func(t *testing.T) {
		t.Run("Serialization", func(t *testing.T) {
			toPtr := utils.ToPointer[PlanetSupportUndefinedWithDefault]
			cfg := utils.TestConfig{SupportUndefined: true, HasDefault: true}
			testCases := []utils.TestCase{
				{From: "PlanetSupportUndefinedWithDefault(9)", Enum: toPtr(9), Expected: utils.Expected{AsSerialized: "PlanetSupportUndefinedWithDefault(9)", IsInvalid: true}},
				{From: "", Enum: toPtr(0), Expected: utils.Expected{AsSerialized: "Earth"}},
				{From: "Earth", Enum: toPtr(0), Expected: utils.Expected{AsSerialized: "Earth"}},
				{From: "Mars", Enum: toPtr(PlanetSupportUndefinedWithDefaultMars), Expected: utils.Expected{AsSerialized: "Mars"}},
				{From: "Pluto", Enum: toPtr(PlanetSupportUndefinedWithDefaultPluto), Expected: utils.Expected{AsSerialized: "Pluto"}},
				{From: "Venus", Enum: toPtr(PlanetSupportUndefinedWithDefaultVenus), Expected: utils.Expected{AsSerialized: "Venus"}},
				{From: "Mercury", Enum: toPtr(PlanetSupportUndefinedWithDefaultMercury), Expected: utils.Expected{AsSerialized: "Mercury"}},
				{From: "Jupiter", Enum: toPtr(PlanetSupportUndefinedWithDefaultJupiter), Expected: utils.Expected{AsSerialized: "Jupiter"}},
				{From: "Saturn", Enum: toPtr(PlanetSupportUndefinedWithDefaultSaturn), Expected: utils.Expected{AsSerialized: "Saturn"}},
				{From: "Uranus", Enum: toPtr(PlanetSupportUndefinedWithDefaultUranus), Expected: utils.Expected{AsSerialized: "Uranus"}},
				{From: "Neptune", Enum: toPtr(PlanetSupportUndefinedWithDefaultNeptune), Expected: utils.Expected{AsSerialized: "Neptune"}},
			}
			for idx, tC := range testCases {
				serializers := []string{"binary", "gql", "json", "sql", "text", "yaml"}
				utils.AssertSerializationInterfacesFor[PlanetSupportUndefinedWithDefault](t, idx, tC, cfg, serializers)
			}
		})
	})
}
