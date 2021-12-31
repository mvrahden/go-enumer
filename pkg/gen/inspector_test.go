package gen

import (
	"fmt"
	"path"
	"path/filepath"
	"testing"

	"github.com/mvrahden/go-enumer/about"
	"github.com/stretchr/testify/require"
)

const (
	packageBase = about.Repo
)

func TestInspectorOutput(t *testing.T) {

	target := "pills"
	pkg := path.Join(packageBase, "/pkg/gen/examples/", target)
	testdatadir := filepath.Join("examples/", target)

	goPkg, err := loadPackage(pkg)
	require.NoError(t, err)

	testCases := []struct {
		NamePrefix string
		Type       GoType
		Offset     uint64
	}{
		{NamePrefix: "Pill", Type: GoTypeSignedInteger},
		{NamePrefix: "PillSigned8", Type: GoTypeSignedInteger8},
		{NamePrefix: "PillSigned16", Type: GoTypeSignedInteger16},
		{NamePrefix: "PillSigned32", Type: GoTypeSignedInteger32},
		{NamePrefix: "PillSigned64", Type: GoTypeSignedInteger64},
		{NamePrefix: "PillUnsigned", Type: GoTypeUnsignedInteger},
		{NamePrefix: "PillUnsigned8", Type: GoTypeUnsignedInteger8},
		{NamePrefix: "PillUnsigned16", Type: GoTypeUnsignedInteger16},
		{NamePrefix: "PillUnsigned32", Type: GoTypeUnsignedInteger32},
		{NamePrefix: "PillUnsigned64", Type: GoTypeUnsignedInteger64},
		{NamePrefix: "PillRowed", Type: GoTypeSignedInteger},
		{NamePrefix: "PillAliased", Type: GoTypeSignedInteger},
		{NamePrefix: "PillUndefined", Type: GoTypeSignedInteger, Offset: 1},
	}

	for idx, tC := range testCases {
		t.Run(fmt.Sprintf("Generate for package \"pills\" (TestCase: %d %q)", idx, tC.NamePrefix), func(t *testing.T) {

			cfg := getConfig(t, testdatadir)
			cfg.TypeAliasName = tC.NamePrefix

			i := NewInspector(cfg)
			f, err := i.Inspect(goPkg)
			require.NoError(t, err)
			require.Len(t, f.ValueSpecs, 4)

			require.Equal(t, &ValueSpec{
				Index:          0,
				IdentifierName: fmt.Sprintf("%sPlacebo", tC.NamePrefix),
				NameString:     "Placebo",
				Type:           tC.Type,
				Value:          0 + tC.Offset,
			}, f.ValueSpecs[0])
			require.Equal(t, &ValueSpec{
				Index:          1,
				IdentifierName: fmt.Sprintf("%sAspirin", tC.NamePrefix),
				NameString:     "Aspirin",
				Type:           tC.Type,
				Value:          1 + tC.Offset,
			}, f.ValueSpecs[1])
			require.Equal(t, &ValueSpec{
				Index:          2,
				IdentifierName: fmt.Sprintf("%sIbuprofen", tC.NamePrefix),
				NameString:     "Ibuprofen",
				Type:           tC.Type,
				Value:          2 + tC.Offset,
			}, f.ValueSpecs[2])
			require.Equal(t, &ValueSpec{
				Index:          3,
				IdentifierName: fmt.Sprintf("%sAcetaminophen", tC.NamePrefix),
				NameString:     "Acetaminophen",
				Type:           tC.Type,
				Value:          3 + tC.Offset,
			}, f.ValueSpecs[3])
		})
	}
}
