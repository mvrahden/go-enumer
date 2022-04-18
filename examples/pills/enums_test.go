package pills

import (
	"fmt"
	"testing"

	"github.com/mvrahden/go-enumer/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestEnums(t *testing.T) {
	t.Run("PillUnsigned", func(t *testing.T) {
		t.Run("Value Sets", func(t *testing.T) {
			require.Equal(t,
				[]string{"PLACEBO", "ASPIRIN", "IBUPROFEN", "PARACETAMOL", "VITAMIN-C"},
				PillUnsignedStrings())
			require.Equal(t,
				[]PillUnsigned{PillUnsignedPlacebo, PillUnsignedAspirin, PillUnsignedIbuprofen, PillUnsignedParacetamol, PillUnsignedVitaminC},
				PillUnsignedValues())
		})
		t.Run("Lookup", func(t *testing.T) {
			type testCase struct {
				enum  PillUnsigned
				upper string
				lower string
			}
			testCases := []testCase{
				{PillUnsignedPlacebo, "PLACEBO", "placebo"},
				{PillUnsignedAspirin, "ASPIRIN", "aspirin"},
				{PillUnsignedIbuprofen, "IBUPROFEN", "ibuprofen"},
				{PillUnsignedParacetamol, "PARACETAMOL", "paracetamol"},
				{PillUnsignedVitaminC, "VITAMIN-C", "vitamin-c"},
			}
			for idx, tC := range testCases {
				t.Run(fmt.Sprintf("Case-sensitive lookup (idx: %d %s)", idx, tC.enum), func(t *testing.T) {
					actual, ok := PillUnsignedFromString(tC.upper)
					require.True(t, ok)
					require.Equal(t, tC.enum, actual)
					actual, ok = PillUnsignedFromString(tC.lower)
					require.False(t, ok)
					require.Equal(t, PillUnsigned(0), actual)
				})
				t.Run(fmt.Sprintf("Case-insensitive lookup (idx: %d %s)", idx, tC.enum), func(t *testing.T) {
					enum, ok := PillUnsignedFromStringIgnoreCase(tC.upper)
					require.True(t, ok)
					require.Equal(t, tC.enum, enum)
					enum, ok = PillUnsignedFromStringIgnoreCase(tC.lower)
					require.True(t, ok)
					require.Equal(t, tC.enum, enum)
				})
			}
		})
		t.Run("Serialization", func(t *testing.T) {
			cfg := utils.TestConfig{}
			toPtr := utils.ToPointer[PillUnsigned]
			testCases := []utils.TestCase{
				{From: "", Enum: toPtr(5), Expected: utils.Expected{AsSerialized: "PillUnsigned(5)", IsInvalid: true}},
				{From: "PLACEBO", Enum: toPtr(0), Expected: utils.Expected{AsSerialized: "PLACEBO"}},
				{From: "ASPIRIN", Enum: toPtr(PillUnsignedAspirin), Expected: utils.Expected{AsSerialized: "ASPIRIN"}},
				{From: "IBUPROFEN", Enum: toPtr(PillUnsignedIbuprofen), Expected: utils.Expected{AsSerialized: "IBUPROFEN"}},
				{From: "PARACETAMOL", Enum: toPtr(PillUnsignedParacetamol), Expected: utils.Expected{AsSerialized: "PARACETAMOL"}},
				{From: "ACETAMINOPHEN", Enum: toPtr(PillUnsignedAcetaminophen), Expected: utils.Expected{AsSerialized: "PARACETAMOL"}},
				{From: "VITAMIN-C", Enum: toPtr(PillUnsignedVitaminC), Expected: utils.Expected{AsSerialized: "VITAMIN-C"}},
			}
			for idx, tC := range testCases {
				serializers := []string{"binary", "gql", "json", "sql", "text", "yaml.v3"}
				utils.AssertSerializationInterfacesFor[PillUnsigned](t, idx, tC, cfg, serializers)
			}
		})
	})
	t.Run("PillAliased", func(t *testing.T) {
		t.Run("Value Sets", func(t *testing.T) {
			require.Equal(t,
				[]string{"PLACEBO", "ASPIRIN", "IBUPROFEN", "PARACETAMOL", "VITAMIN-C"},
				PillAliasedStrings())
			require.Equal(t,
				[]PillAliased{PillAliasedPlacebo, PillAliasedAspirin, PillAliasedIbuprofen, PillAliasedParacetamol, PillAliasedVitaminC},
				PillAliasedValues())
		})
		t.Run("Lookup", func(t *testing.T) {
			type testCase struct {
				enum  PillAliased
				upper string
				lower string
			}
			testCases := []testCase{
				{PillAliasedPlacebo, "PLACEBO", "placebo"},
				{PillAliasedAspirin, "ASPIRIN", "aspirin"},
				{PillAliasedIbuprofen, "IBUPROFEN", "ibuprofen"},
				{PillAliasedParacetamol, "PARACETAMOL", "paracetamol"},
				{PillAliasedVitaminC, "VITAMIN-C", "vitamin-c"},
			}
			for idx, tC := range testCases {
				t.Run(fmt.Sprintf("Case-sensitive lookup (idx: %d %s)", idx, tC.enum), func(t *testing.T) {
					actual, ok := PillAliasedFromString(tC.upper)
					require.True(t, ok)
					require.Equal(t, tC.enum, actual)
					actual, ok = PillAliasedFromString(tC.lower)
					require.False(t, ok)
					require.Equal(t, PillAliased(0), actual)
				})
				t.Run(fmt.Sprintf("Case-insensitive lookup (idx: %d %s)", idx, tC.enum), func(t *testing.T) {
					enum, ok := PillAliasedFromStringIgnoreCase(tC.upper)
					require.True(t, ok)
					require.Equal(t, tC.enum, enum)
					enum, ok = PillAliasedFromStringIgnoreCase(tC.lower)
					require.True(t, ok)
					require.Equal(t, tC.enum, enum)
				})
			}
		})
		t.Run("Serialization", func(t *testing.T) {
			cfg := utils.TestConfig{}
			toPtr := utils.ToPointer[PillAliased]
			testCases := []utils.TestCase{
				{From: "", Enum: toPtr(5), Expected: utils.Expected{AsSerialized: "PillAliased(5)", IsInvalid: true}},
				{From: "PLACEBO", Enum: toPtr(0), Expected: utils.Expected{AsSerialized: "PLACEBO"}},
				{From: "ASPIRIN", Enum: toPtr(PillAliasedAspirin), Expected: utils.Expected{AsSerialized: "ASPIRIN"}},
				{From: "IBUPROFEN", Enum: toPtr(PillAliasedIbuprofen), Expected: utils.Expected{AsSerialized: "IBUPROFEN"}},
				{From: "PARACETAMOL", Enum: toPtr(PillAliasedParacetamol), Expected: utils.Expected{AsSerialized: "PARACETAMOL"}},
				{From: "ACETAMINOPHEN", Enum: toPtr(PillAliasedAcetaminophen), Expected: utils.Expected{AsSerialized: "PARACETAMOL"}},
				{From: "VITAMIN-C", Enum: toPtr(PillAliasedVitaminC), Expected: utils.Expected{AsSerialized: "VITAMIN-C"}},
			}
			for idx, tC := range testCases {
				serializers := []string{"binary", "gql", "json", "sql", "text", "yaml.v3"}
				utils.AssertSerializationInterfacesFor[PillAliased](t, idx, tC, cfg, serializers)
			}
		})
	})
	t.Run("PillRowed", func(t *testing.T) {
		t.Run("Value Sets", func(t *testing.T) {
			require.Equal(t,
				[]string{"PLACEBO", "ASPIRIN", "IBUPROFEN", "PARACETAMOL", "VITAMIN-C"},
				PillRowedStrings())
			require.Equal(t,
				[]PillRowed{PillRowedPlacebo, PillRowedAspirin, PillRowedIbuprofen, PillRowedParacetamol, PillRowedVitaminC},
				PillRowedValues())
		})
		t.Run("Lookup", func(t *testing.T) {
			type testCase struct {
				enum  PillRowed
				upper string
				lower string
			}
			testCases := []testCase{
				{PillRowedPlacebo, "PLACEBO", "placebo"},
				{PillRowedAspirin, "ASPIRIN", "aspirin"},
				{PillRowedIbuprofen, "IBUPROFEN", "ibuprofen"},
				{PillRowedParacetamol, "PARACETAMOL", "paracetamol"},
				{PillRowedVitaminC, "VITAMIN-C", "vitamin-c"},
			}
			for idx, tC := range testCases {
				t.Run(fmt.Sprintf("Case-sensitive lookup (idx: %d %s)", idx, tC.enum), func(t *testing.T) {
					actual, ok := PillRowedFromString(tC.upper)
					require.True(t, ok)
					require.Equal(t, tC.enum, actual)
					actual, ok = PillRowedFromString(tC.lower)
					require.False(t, ok)
					require.Equal(t, PillRowed(0), actual)
				})
				t.Run(fmt.Sprintf("Case-insensitive lookup (idx: %d %s)", idx, tC.enum), func(t *testing.T) {
					enum, ok := PillRowedFromStringIgnoreCase(tC.upper)
					require.True(t, ok)
					require.Equal(t, tC.enum, enum)
					enum, ok = PillRowedFromStringIgnoreCase(tC.lower)
					require.True(t, ok)
					require.Equal(t, tC.enum, enum)
				})
			}
		})
		t.Run("Serialization", func(t *testing.T) {
			cfg := utils.TestConfig{}
			toPtr := utils.ToPointer[PillRowed]
			testCases := []utils.TestCase{
				{From: "", Enum: toPtr(5), Expected: utils.Expected{AsSerialized: "PillRowed(5)", IsInvalid: true}},
				{From: "PLACEBO", Enum: toPtr(0), Expected: utils.Expected{AsSerialized: "PLACEBO"}},
				{From: "ASPIRIN", Enum: toPtr(PillRowedAspirin), Expected: utils.Expected{AsSerialized: "ASPIRIN"}},
				{From: "IBUPROFEN", Enum: toPtr(PillRowedIbuprofen), Expected: utils.Expected{AsSerialized: "IBUPROFEN"}},
				{From: "PARACETAMOL", Enum: toPtr(PillRowedParacetamol), Expected: utils.Expected{AsSerialized: "PARACETAMOL"}},
				{From: "ACETAMINOPHEN", Enum: toPtr(PillRowedAcetaminophen), Expected: utils.Expected{AsSerialized: "PARACETAMOL"}},
				{From: "VITAMIN-C", Enum: toPtr(PillRowedVitaminC), Expected: utils.Expected{AsSerialized: "VITAMIN-C"}},
			}
			for idx, tC := range testCases {
				serializers := []string{"binary", "gql", "json", "sql", "text", "yaml.v3"}
				utils.AssertSerializationInterfacesFor[PillRowed](t, idx, tC, cfg, serializers)
			}
		})
	})
	t.Run("PillUnsigned8", func(t *testing.T) {
		t.Run("Value Sets", func(t *testing.T) {
			require.Equal(t,
				[]string{"PLACEBO", "ASPIRIN", "IBUPROFEN", "PARACETAMOL", "VITAMIN-C"},
				PillUnsigned8Strings())
			require.Equal(t,
				[]PillUnsigned8{PillUnsigned8Placebo, PillUnsigned8Aspirin, PillUnsigned8Ibuprofen, PillUnsigned8Paracetamol, PillUnsigned8VitaminC},
				PillUnsigned8Values())
		})
		t.Run("Lookup", func(t *testing.T) {
			type testCase struct {
				enum  PillUnsigned8
				upper string
				lower string
			}
			testCases := []testCase{
				{PillUnsigned8Placebo, "PLACEBO", "placebo"},
				{PillUnsigned8Aspirin, "ASPIRIN", "aspirin"},
				{PillUnsigned8Ibuprofen, "IBUPROFEN", "ibuprofen"},
				{PillUnsigned8Paracetamol, "PARACETAMOL", "paracetamol"},
				{PillUnsigned8VitaminC, "VITAMIN-C", "vitamin-c"},
			}
			for idx, tC := range testCases {
				t.Run(fmt.Sprintf("Case-sensitive lookup (idx: %d %s)", idx, tC.enum), func(t *testing.T) {
					actual, ok := PillUnsigned8FromString(tC.upper)
					require.True(t, ok)
					require.Equal(t, tC.enum, actual)
					actual, ok = PillUnsigned8FromString(tC.lower)
					require.False(t, ok)
					require.Equal(t, PillUnsigned8(0), actual)
				})
				t.Run(fmt.Sprintf("Case-insensitive lookup (idx: %d %s)", idx, tC.enum), func(t *testing.T) {
					enum, ok := PillUnsigned8FromStringIgnoreCase(tC.upper)
					require.True(t, ok)
					require.Equal(t, tC.enum, enum)
					enum, ok = PillUnsigned8FromStringIgnoreCase(tC.lower)
					require.True(t, ok)
					require.Equal(t, tC.enum, enum)
				})
			}
		})
		t.Run("Serialization", func(t *testing.T) {
			cfg := utils.TestConfig{}
			toPtr := utils.ToPointer[PillUnsigned8]
			testCases := []utils.TestCase{
				{From: "", Enum: toPtr(5), Expected: utils.Expected{AsSerialized: "PillUnsigned8(5)", IsInvalid: true}},
				{From: "PLACEBO", Enum: toPtr(0), Expected: utils.Expected{AsSerialized: "PLACEBO"}},
				{From: "ASPIRIN", Enum: toPtr(PillUnsigned8Aspirin), Expected: utils.Expected{AsSerialized: "ASPIRIN"}},
				{From: "IBUPROFEN", Enum: toPtr(PillUnsigned8Ibuprofen), Expected: utils.Expected{AsSerialized: "IBUPROFEN"}},
				{From: "PARACETAMOL", Enum: toPtr(PillUnsigned8Paracetamol), Expected: utils.Expected{AsSerialized: "PARACETAMOL"}},
				{From: "ACETAMINOPHEN", Enum: toPtr(PillUnsigned8Acetaminophen), Expected: utils.Expected{AsSerialized: "PARACETAMOL"}},
				{From: "VITAMIN-C", Enum: toPtr(PillUnsigned8VitaminC), Expected: utils.Expected{AsSerialized: "VITAMIN-C"}},
			}
			for idx, tC := range testCases {
				serializers := []string{"binary", "gql", "json", "sql", "text", "yaml.v3"}
				utils.AssertSerializationInterfacesFor[PillUnsigned8](t, idx, tC, cfg, serializers)
			}
		})
	})
	t.Run("PillUnsigned16", func(t *testing.T) {
		t.Run("Value Sets", func(t *testing.T) {
			require.Equal(t,
				[]string{"PLACEBO", "ASPIRIN", "IBUPROFEN", "PARACETAMOL", "VITAMIN-C"},
				PillUnsigned16Strings())
			require.Equal(t,
				[]PillUnsigned16{PillUnsigned16Placebo, PillUnsigned16Aspirin, PillUnsigned16Ibuprofen, PillUnsigned16Paracetamol, PillUnsigned16VitaminC},
				PillUnsigned16Values())
		})
		t.Run("Lookup", func(t *testing.T) {
			type testCase struct {
				enum  PillUnsigned16
				upper string
				lower string
			}
			testCases := []testCase{
				{PillUnsigned16Placebo, "PLACEBO", "placebo"},
				{PillUnsigned16Aspirin, "ASPIRIN", "aspirin"},
				{PillUnsigned16Ibuprofen, "IBUPROFEN", "ibuprofen"},
				{PillUnsigned16Paracetamol, "PARACETAMOL", "paracetamol"},
				{PillUnsigned16VitaminC, "VITAMIN-C", "vitamin-c"},
			}
			for idx, tC := range testCases {
				t.Run(fmt.Sprintf("Case-sensitive lookup (idx: %d %s)", idx, tC.enum), func(t *testing.T) {
					actual, ok := PillUnsigned16FromString(tC.upper)
					require.True(t, ok)
					require.Equal(t, tC.enum, actual)
					actual, ok = PillUnsigned16FromString(tC.lower)
					require.False(t, ok)
					require.Equal(t, PillUnsigned16(0), actual)
				})
				t.Run(fmt.Sprintf("Case-insensitive lookup (idx: %d %s)", idx, tC.enum), func(t *testing.T) {
					enum, ok := PillUnsigned16FromStringIgnoreCase(tC.upper)
					require.True(t, ok)
					require.Equal(t, tC.enum, enum)
					enum, ok = PillUnsigned16FromStringIgnoreCase(tC.lower)
					require.True(t, ok)
					require.Equal(t, tC.enum, enum)
				})
			}
		})
		t.Run("Serialization", func(t *testing.T) {
			cfg := utils.TestConfig{}
			toPtr := utils.ToPointer[PillUnsigned16]
			testCases := []utils.TestCase{
				{From: "", Enum: toPtr(5), Expected: utils.Expected{AsSerialized: "PillUnsigned16(5)", IsInvalid: true}},
				{From: "PLACEBO", Enum: toPtr(0), Expected: utils.Expected{AsSerialized: "PLACEBO"}},
				{From: "ASPIRIN", Enum: toPtr(PillUnsigned16Aspirin), Expected: utils.Expected{AsSerialized: "ASPIRIN"}},
				{From: "IBUPROFEN", Enum: toPtr(PillUnsigned16Ibuprofen), Expected: utils.Expected{AsSerialized: "IBUPROFEN"}},
				{From: "PARACETAMOL", Enum: toPtr(PillUnsigned16Paracetamol), Expected: utils.Expected{AsSerialized: "PARACETAMOL"}},
				{From: "ACETAMINOPHEN", Enum: toPtr(PillUnsigned16Acetaminophen), Expected: utils.Expected{AsSerialized: "PARACETAMOL"}},
				{From: "VITAMIN-C", Enum: toPtr(PillUnsigned16VitaminC), Expected: utils.Expected{AsSerialized: "VITAMIN-C"}},
			}
			for idx, tC := range testCases {
				serializers := []string{"binary", "gql", "json", "sql", "text", "yaml.v3"}
				utils.AssertSerializationInterfacesFor[PillUnsigned16](t, idx, tC, cfg, serializers)
			}
		})
	})
	t.Run("PillUnsigned32", func(t *testing.T) {
		t.Run("Value Sets", func(t *testing.T) {
			require.Equal(t,
				[]string{"PLACEBO", "ASPIRIN", "IBUPROFEN", "PARACETAMOL", "VITAMIN-C"},
				PillUnsigned32Strings())
			require.Equal(t,
				[]PillUnsigned32{PillUnsigned32Placebo, PillUnsigned32Aspirin, PillUnsigned32Ibuprofen, PillUnsigned32Paracetamol, PillUnsigned32VitaminC},
				PillUnsigned32Values())
		})
		t.Run("Lookup", func(t *testing.T) {
			type testCase struct {
				enum  PillUnsigned32
				upper string
				lower string
			}
			testCases := []testCase{
				{PillUnsigned32Placebo, "PLACEBO", "placebo"},
				{PillUnsigned32Aspirin, "ASPIRIN", "aspirin"},
				{PillUnsigned32Ibuprofen, "IBUPROFEN", "ibuprofen"},
				{PillUnsigned32Paracetamol, "PARACETAMOL", "paracetamol"},
				{PillUnsigned32VitaminC, "VITAMIN-C", "vitamin-c"},
			}
			for idx, tC := range testCases {
				t.Run(fmt.Sprintf("Case-sensitive lookup (idx: %d %s)", idx, tC.enum), func(t *testing.T) {
					actual, ok := PillUnsigned32FromString(tC.upper)
					require.True(t, ok)
					require.Equal(t, tC.enum, actual)
					actual, ok = PillUnsigned32FromString(tC.lower)
					require.False(t, ok)
					require.Equal(t, PillUnsigned32(0), actual)
				})
				t.Run(fmt.Sprintf("Case-insensitive lookup (idx: %d %s)", idx, tC.enum), func(t *testing.T) {
					enum, ok := PillUnsigned32FromStringIgnoreCase(tC.upper)
					require.True(t, ok)
					require.Equal(t, tC.enum, enum)
					enum, ok = PillUnsigned32FromStringIgnoreCase(tC.lower)
					require.True(t, ok)
					require.Equal(t, tC.enum, enum)
				})
			}
		})
		t.Run("Serialization", func(t *testing.T) {
			cfg := utils.TestConfig{}
			toPtr := utils.ToPointer[PillUnsigned32]
			testCases := []utils.TestCase{
				{From: "", Enum: toPtr(5), Expected: utils.Expected{AsSerialized: "PillUnsigned32(5)", IsInvalid: true}},
				{From: "PLACEBO", Enum: toPtr(0), Expected: utils.Expected{AsSerialized: "PLACEBO"}},
				{From: "ASPIRIN", Enum: toPtr(PillUnsigned32Aspirin), Expected: utils.Expected{AsSerialized: "ASPIRIN"}},
				{From: "IBUPROFEN", Enum: toPtr(PillUnsigned32Ibuprofen), Expected: utils.Expected{AsSerialized: "IBUPROFEN"}},
				{From: "PARACETAMOL", Enum: toPtr(PillUnsigned32Paracetamol), Expected: utils.Expected{AsSerialized: "PARACETAMOL"}},
				{From: "ACETAMINOPHEN", Enum: toPtr(PillUnsigned32Acetaminophen), Expected: utils.Expected{AsSerialized: "PARACETAMOL"}},
				{From: "VITAMIN-C", Enum: toPtr(PillUnsigned32VitaminC), Expected: utils.Expected{AsSerialized: "VITAMIN-C"}},
			}
			for idx, tC := range testCases {
				serializers := []string{"binary", "gql", "json", "sql", "text", "yaml.v3"}
				utils.AssertSerializationInterfacesFor[PillUnsigned32](t, idx, tC, cfg, serializers)
			}
		})
	})
	t.Run("PillUnsigned64", func(t *testing.T) {
		t.Run("Value Sets", func(t *testing.T) {
			require.Equal(t,
				[]string{"PLACEBO", "ASPIRIN", "IBUPROFEN", "PARACETAMOL", "VITAMIN-C"},
				PillUnsigned64Strings())
			require.Equal(t,
				[]PillUnsigned64{PillUnsigned64Placebo, PillUnsigned64Aspirin, PillUnsigned64Ibuprofen, PillUnsigned64Paracetamol, PillUnsigned64VitaminC},
				PillUnsigned64Values())
		})
		t.Run("Lookup", func(t *testing.T) {
			type testCase struct {
				enum  PillUnsigned64
				upper string
				lower string
			}
			testCases := []testCase{
				{PillUnsigned64Placebo, "PLACEBO", "placebo"},
				{PillUnsigned64Aspirin, "ASPIRIN", "aspirin"},
				{PillUnsigned64Ibuprofen, "IBUPROFEN", "ibuprofen"},
				{PillUnsigned64Paracetamol, "PARACETAMOL", "paracetamol"},
				{PillUnsigned64VitaminC, "VITAMIN-C", "vitamin-c"},
			}
			for idx, tC := range testCases {
				t.Run(fmt.Sprintf("Case-sensitive lookup (idx: %d %s)", idx, tC.enum), func(t *testing.T) {
					actual, ok := PillUnsigned64FromString(tC.upper)
					require.True(t, ok)
					require.Equal(t, tC.enum, actual)
					actual, ok = PillUnsigned64FromString(tC.lower)
					require.False(t, ok)
					require.Equal(t, PillUnsigned64(0), actual)
				})
				t.Run(fmt.Sprintf("Case-insensitive lookup (idx: %d %s)", idx, tC.enum), func(t *testing.T) {
					enum, ok := PillUnsigned64FromStringIgnoreCase(tC.upper)
					require.True(t, ok)
					require.Equal(t, tC.enum, enum)
					enum, ok = PillUnsigned64FromStringIgnoreCase(tC.lower)
					require.True(t, ok)
					require.Equal(t, tC.enum, enum)
				})
			}
		})
		t.Run("Serialization", func(t *testing.T) {
			cfg := utils.TestConfig{}
			toPtr := utils.ToPointer[PillUnsigned64]
			testCases := []utils.TestCase{
				{From: "", Enum: toPtr(5), Expected: utils.Expected{AsSerialized: "PillUnsigned64(5)", IsInvalid: true}},
				{From: "PLACEBO", Enum: toPtr(0), Expected: utils.Expected{AsSerialized: "PLACEBO"}},
				{From: "ASPIRIN", Enum: toPtr(PillUnsigned64Aspirin), Expected: utils.Expected{AsSerialized: "ASPIRIN"}},
				{From: "IBUPROFEN", Enum: toPtr(PillUnsigned64Ibuprofen), Expected: utils.Expected{AsSerialized: "IBUPROFEN"}},
				{From: "PARACETAMOL", Enum: toPtr(PillUnsigned64Paracetamol), Expected: utils.Expected{AsSerialized: "PARACETAMOL"}},
				{From: "ACETAMINOPHEN", Enum: toPtr(PillUnsigned64Acetaminophen), Expected: utils.Expected{AsSerialized: "PARACETAMOL"}},
				{From: "VITAMIN-C", Enum: toPtr(PillUnsigned64VitaminC), Expected: utils.Expected{AsSerialized: "VITAMIN-C"}},
			}
			for idx, tC := range testCases {
				serializers := []string{"binary", "gql", "json", "sql", "text", "yaml.v3"}
				utils.AssertSerializationInterfacesFor[PillUnsigned64](t, idx, tC, cfg, serializers)
			}
		})
	})
}
