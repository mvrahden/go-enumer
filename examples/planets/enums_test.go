package planets

import (
	"fmt"
	"testing"

	"github.com/mvrahden/go-enumer/pkg/utils"
	"github.com/stretchr/testify/require"
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
				t.Run(fmt.Sprintf("Serializers (idx: %d from %q)", idx, tC.From), func(t *testing.T) {
					utils.AssertSerializers[Planet](t, tC, "binary")
					utils.AssertSerializers[Planet](t, tC, "gql")
					utils.AssertSerializers[Planet](t, tC, "json")
					utils.AssertSerializers[Planet](t, tC, "sql")
					utils.AssertSerializers[Planet](t, tC, "text")
					utils.AssertSerializers[Planet](t, tC, "yaml")
				})
				t.Run(fmt.Sprintf("Deserializers (idx: %d from %q)", idx, tC.From), func(t *testing.T) {
					utils.AssertDeserializers(t, tC, cfg, "binary", utils.ZeroValuer[Planet])
					utils.AssertDeserializers(t, tC, cfg, "gql", utils.ZeroValuer[Planet])
					utils.AssertDeserializers(t, tC, cfg, "json", utils.ZeroValuer[Planet])
					utils.AssertDeserializers(t, tC, cfg, "sql", utils.ZeroValuer[Planet])
					utils.AssertDeserializers(t, tC, cfg, "text", utils.ZeroValuer[Planet])
					utils.AssertDeserializers(t, tC, cfg, "yaml", utils.ZeroValuer[Planet])
				})
			}
		})
	})
	t.Run("PlanetWithDefault", func(t *testing.T) {
		t.Run("Serialization", func(t *testing.T) {
			cfg := utils.TestConfig{}
			toPtr := utils.ToPointer[PlanetWithDefault]
			testCases := []utils.TestCase{
				{From: "PlanetWithDefault(9)", Enum: toPtr(9), Expected: utils.Expected{AsSerialized: "PlanetWithDefault(9)", IsInvalid: true}},
				// TODO: fix these tests
				// {From: "", Enum: toPtr(0), Expected: utils.Expected{AsSerialized: "PlanetWithDefault(0)", invalid: true}},
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
				t.Run(fmt.Sprintf("Serializers (idx: %d from %q)", idx, tC.From), func(t *testing.T) {
					utils.AssertSerializers[PlanetWithDefault](t, tC, "binary")
					utils.AssertSerializers[PlanetWithDefault](t, tC, "gql")
					utils.AssertSerializers[PlanetWithDefault](t, tC, "json")
					utils.AssertSerializers[PlanetWithDefault](t, tC, "sql")
					utils.AssertSerializers[PlanetWithDefault](t, tC, "text")
					utils.AssertSerializers[PlanetWithDefault](t, tC, "yaml")
				})
				t.Run(fmt.Sprintf("Deserializers (idx: %d from %q)", idx, tC.From), func(t *testing.T) {
					utils.AssertDeserializers(t, tC, cfg, "binary", utils.ZeroValuer[PlanetWithDefault])
					utils.AssertDeserializers(t, tC, cfg, "gql", utils.ZeroValuer[PlanetWithDefault])
					utils.AssertDeserializers(t, tC, cfg, "json", utils.ZeroValuer[PlanetWithDefault])
					utils.AssertDeserializers(t, tC, cfg, "sql", utils.ZeroValuer[PlanetWithDefault])
					utils.AssertDeserializers(t, tC, cfg, "text", utils.ZeroValuer[PlanetWithDefault])
					utils.AssertDeserializers(t, tC, cfg, "yaml", utils.ZeroValuer[PlanetWithDefault])
				})
			}
		})
	})
	t.Run("PlanetSupportUndefined", func(t *testing.T) {
		t.Run(`generates "<type>Undefined" constant`, func(t *testing.T) {
			// asserts that special const "<type>Undefined" is generated.
			require.Equal(t, 0, int(PlanetSupportUndefinedUndefined))
		})
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
				t.Run(fmt.Sprintf("Serializers (idx: %d from %q)", idx, tC.From), func(t *testing.T) {
					utils.AssertSerializers[PlanetSupportUndefined](t, tC, "binary")
					utils.AssertSerializers[PlanetSupportUndefined](t, tC, "gql")
					utils.AssertSerializers[PlanetSupportUndefined](t, tC, "json")
					utils.AssertSerializers[PlanetSupportUndefined](t, tC, "sql")
					utils.AssertSerializers[PlanetSupportUndefined](t, tC, "text")
					utils.AssertSerializers[PlanetSupportUndefined](t, tC, "yaml")
				})
				t.Run(fmt.Sprintf("Deserializers (idx: %d from %q)", idx, tC.From), func(t *testing.T) {
					utils.AssertDeserializers(t, tC, cfg, "binary", utils.ZeroValuer[PlanetSupportUndefined])
					utils.AssertDeserializers(t, tC, cfg, "gql", utils.ZeroValuer[PlanetSupportUndefined])
					utils.AssertDeserializers(t, tC, cfg, "json", utils.ZeroValuer[PlanetSupportUndefined])
					utils.AssertDeserializers(t, tC, cfg, "sql", utils.ZeroValuer[PlanetSupportUndefined])
					utils.AssertDeserializers(t, tC, cfg, "text", utils.ZeroValuer[PlanetSupportUndefined])
					utils.AssertDeserializers(t, tC, cfg, "yaml", utils.ZeroValuer[PlanetSupportUndefined])
				})
			}
		})
	})
	t.Run("PlanetSupportUndefinedWithDefault", func(t *testing.T) {
		t.Run("Serialization", func(t *testing.T) {
			toPtr := utils.ToPointer[PlanetSupportUndefinedWithDefault]
			cfg := utils.TestConfig{SupportUndefined: true}
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
				t.Run(fmt.Sprintf("Serializers (idx: %d from %q)", idx, tC.From), func(t *testing.T) {
					utils.AssertSerializers[PlanetSupportUndefinedWithDefault](t, tC, "binary")
					utils.AssertSerializers[PlanetSupportUndefinedWithDefault](t, tC, "gql")
					utils.AssertSerializers[PlanetSupportUndefinedWithDefault](t, tC, "json")
					utils.AssertSerializers[PlanetSupportUndefinedWithDefault](t, tC, "sql")
					utils.AssertSerializers[PlanetSupportUndefinedWithDefault](t, tC, "text")
					utils.AssertSerializers[PlanetSupportUndefinedWithDefault](t, tC, "yaml")
				})
				t.Run(fmt.Sprintf("Deserializers (idx: %d from %q)", idx, tC.From), func(t *testing.T) {
					utils.AssertDeserializers(t, tC, cfg, "binary", utils.ZeroValuer[PlanetSupportUndefinedWithDefault])
					utils.AssertDeserializers(t, tC, cfg, "gql", utils.ZeroValuer[PlanetSupportUndefinedWithDefault])
					utils.AssertDeserializers(t, tC, cfg, "json", utils.ZeroValuer[PlanetSupportUndefinedWithDefault])
					utils.AssertDeserializers(t, tC, cfg, "sql", utils.ZeroValuer[PlanetSupportUndefinedWithDefault])
					utils.AssertDeserializers(t, tC, cfg, "text", utils.ZeroValuer[PlanetSupportUndefinedWithDefault])
					utils.AssertDeserializers(t, tC, cfg, "yaml", utils.ZeroValuer[PlanetSupportUndefinedWithDefault])
				})
			}
		})
	})
}
