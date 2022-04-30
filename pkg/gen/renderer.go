package gen

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/ettle/strcase"
	"github.com/mvrahden/go-enumer/about"
	"github.com/mvrahden/go-enumer/config"
	"github.com/mvrahden/go-enumer/pkg/utils/slices"
)

type renderer struct{}

func NewRenderer(cfg *config.Options) *renderer {
	r := renderer{}
	return &r
}

//go:embed static
var templates embed.FS

var (
	headerTpl = template.Must(template.New("header").ParseFS(templates, "static/header.*"))
	enumTpl   = template.Must(template.New("enum").Funcs(tplFuncs).ParseFS(templates, "static/enum.*"))
)

func (r *renderer) Render(f *File) ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := r.renderFileHeader(buf, f); err != nil {
		return nil, fmt.Errorf("failed rendering file header. err: %w", err)
	}

	idx, err := slices.RangeErr(f.TypeSpecs, func(ts *TypeSpec, _ int) error {
		return r.renderForTypeSpec(buf, ts)
	})
	if err != nil {
		return nil, fmt.Errorf("failed rendering sources for %q. err: %w", f.TypeSpecs[idx].Name, err)
	}
	return buf.Bytes(), nil
}

func (r *renderer) renderFileHeader(buf *bytes.Buffer, f *File) error {
	type TplData struct {
		RepoName          string
		PackageName       string
		Imports           []*Import
		ContainsErrorsPkg bool
	}
	data := TplData{
		RepoName:    about.ShortInfo(),
		PackageName: f.Header.Package.Name,
		Imports:     f.Imports,
		ContainsErrorsPkg: slices.Any(f.Imports, func(v *Import) bool {
			return v.Path == "errors"
		}),
	}
	return headerTpl.ExecuteTemplate(buf, "header.go.tpl", map[string]any{"Header": data})
}

func (r *renderer) renderForTypeSpec(buf *bytes.Buffer, ts *TypeSpec) error {
	if ts.IsFromCsvSource {
		ts.Config.TransformStrategy = "noop"
	}

	util := newRenderUtil(ts.Config)

	for _, v := range ts.ValueSpecs {
		v.EnumValue = util.transform(v.EnumValue)
	}

	requiresGeneratedUndefinedValue := ts.Config.SupportedFeatures.Contains(config.SupportUndefined) &&
		slices.None(ts.ValueSpecs, func(v *ValueSpec) bool { return v.Value == 0 })

	type AdditionalDataHeader struct {
		Name string // hint: the column name as-is (from CSV)
		Type GoType // hint: the type infered by type syntax
	}
	type AdditionalDataCell struct {
		Value string // hint: source representation of the value, e.g. literal strings are quoted
	}
	type AdditionalData struct {
		Headers []AdditionalDataHeader
		Rows    [][]AdditionalDataCell
	}
	type EnumValue struct {
		Value              uint64 // hint: the enum's numeric representation
		String             string // hint: the enum's string representation
		ConstName          string // hint: the enum's constant name
		Position           int    // hint: start index of enum value string within enum aggregate string
		Length             int    // hint: length of
		IsAlternativeValue bool   // hint: is the enum an alternative value
	}
	type Enum struct {
		Name                            string
		Values                          []EnumValue
		RequiresGeneratedUndefinedValue bool
		IsFromCsvSource                 bool
		HasAdditionalData               bool
		AdditionalData                  AdditionalData
	}
	enum := Enum{
		Name: ts.Name,
		Values: slices.Map(ts.ValueSpecs, func(v *ValueSpec, idx int) EnumValue {
			var isAlternativeValue bool
			if idx > 0 {
				isAlternativeValue = v.Value == ts.ValueSpecs[idx-1].Value
			}
			return EnumValue{
				Value:     v.Value,
				String:    v.EnumValue,
				ConstName: v.IdentifierName,
				Position: slices.Reduce(ts.ValueSpecs[0:idx], func(v *ValueSpec, acc int) int {
					return acc + len(v.EnumValue)
				}),
				Length:             len(v.EnumValue),
				IsAlternativeValue: isAlternativeValue,
			}
		}),
		RequiresGeneratedUndefinedValue: requiresGeneratedUndefinedValue,
		IsFromCsvSource:                 ts.IsFromCsvSource,
		HasAdditionalData:               ts.HasAdditionalData,
		AdditionalData: AdditionalData{
			Headers: slices.Map(ts.DataColumns, func(v DataHeader, idx int) AdditionalDataHeader {
				return AdditionalDataHeader{v.Name, v.Type}
			}),
			Rows: slices.Map(ts.ValueSpecs, func(v *ValueSpec, rowIdx int) []AdditionalDataCell {
				return slices.Map(v.DataCells, func(v DataCell, colIdx int) AdditionalDataCell {
					return AdditionalDataCell{ts.DataColumns[colIdx].Type.ToSource(v.ValueString)}
				})
			}),
		},
	}

	{ // write consts
		type TplData struct {
			Enum
			AggregatedValueStrings string
		}

		data := TplData{
			Enum: enum,
			AggregatedValueStrings: slices.ReduceSeed(ts.ValueSpecs, &bytes.Buffer{}, func(v *ValueSpec, acc *bytes.Buffer) *bytes.Buffer {
				acc.WriteString(v.EnumValue)
				return acc
			}).String(),
		}
		if err := enumTpl.ExecuteTemplate(buf, "enum.consts.go.tpl", map[string]any{"Type": data}); err != nil {
			return err
		}
	}

	{ // write vars
		type Extent struct {
			Min uint64 // hint: the lower numerical bound of the enum set
			Max uint64 // hint: the upper numerical bound of the enum set
		}
		type TplData struct {
			Enum
			CountUniqueValues int    // hint: count of all enums, less the alternative values
			Extent            Extent // hint: extent/range of the enum set [min,max]
		}

		lowerBound := ts.ValueSpecs[0].Value
		if requiresGeneratedUndefinedValue {
			lowerBound = 0
		}
		data := TplData{
			Enum: enum,
			CountUniqueValues: slices.Count(ts.ValueSpecs, func(v *ValueSpec, idx int) bool {
				if idx == 0 {
					return true
				}
				return v.Value != ts.ValueSpecs[idx-1].Value
			}),
			Extent: Extent{
				Min: lowerBound,
				Max: ts.ValueSpecs[len(ts.ValueSpecs)-1].Value,
			},
		}
		if err := enumTpl.ExecuteTemplate(buf, "enum.vars.go.tpl", map[string]any{"Type": data}); err != nil {
			return err
		}
	}

	{ // compiletime assertion
		if err := enumTpl.ExecuteTemplate(buf, "enum.assertions.go.tpl", map[string]any{"Type": enum}); err != nil {
			return err
		}
	}

	{ // standard functions
		type TplData struct {
			Enum
			RequiresOffset bool
		}
		data := TplData{
			Enum:           enum,
			RequiresOffset: ts.ValueSpecs[0].Value > 0,
		}

		if err := enumTpl.ExecuteTemplate(buf, "enum.base-funcs.go.tpl", map[string]any{"Type": data}); err != nil {
			return err
		}
	}

	{ // value mappings and lookup
		type TplData struct {
			Enum
			SupportUndefined bool
		}
		data := TplData{
			Enum:             enum,
			SupportUndefined: ts.Config.SupportedFeatures.Contains(config.SupportUndefined),
		}

		if err := enumTpl.ExecuteTemplate(buf, "enum.lookup-funcs.go.tpl", map[string]any{"Type": data}); err != nil {
			return err
		}
	}

	if err := util.renderSerializers(buf, ts); err != nil {
		return err
	}
	return util.renderEntInterfaceSupport(buf, ts)
}

var tplFuncs = template.FuncMap{
	"add": func(a, b int) int {
		return a + b
	},
	"sub": func(a, b int) int {
		return a - b
	},
	"lower":  lowerCaseTransformer,
	"pascal": pascalCaseTransformer,
	"receiver": func(s string) string {
		for _, v := range s {
			return "_" + strings.ToLower(string(v)) // take first rune
		}
		return "_x" // fallback
	},
	"zero": func(t GoType) string {
		return t.ZeroValueString()
	},
	"contains": func(s []string, t string) bool {
		for _, v := range s {
			if v == t {
				return true
			}
		}
		return false
	},
}

func newRenderUtil(cfg *config.Options) *renderUtil {
	r := renderUtil{}
	for _, fn := range []renderUtilOpt{withDefaults, withTransformStrategy} {
		fn(&r, cfg)
	}
	return &r
}

type renderUtilOpt func(r *renderUtil, c *config.Options)

func withDefaults(r *renderUtil, c *config.Options) {
	r.transform = noopCaseTransformer
}

func withTransformStrategy(r *renderUtil, c *config.Options) {
	r.transform = getTransformStrategy(c)
}

type renderUtil struct {
	transform stringCaseTransformer
}

type stringCaseTransformer func(v string) string

var (
	noopCaseTransformer = func(value string) string {
		return value
	}
	pascalCaseTransformer = func(value string) string {
		return strcase.ToPascal(value)
	}
	camelCaseTransformer = func(value string) string {
		return strcase.ToCamel(value)
	}
	kebabCaseTransformer = func(value string) string {
		return strcase.ToKebab(value)
	}
	snakeCaseTransformer = func(value string) string {
		return strcase.ToSnake(value)
	}
	lowerCaseTransformer = func(value string) string {
		return strings.ToLower(value)
	}
	upperCaseTransformer = func(value string) string {
		return strings.ToUpper(value)
	}
	upperSnakeCaseTransformer = func(value string) string {
		return strcase.ToSNAKE(value)
	}
	upperKebabCaseTransformer = func(value string) string {
		return strcase.ToKEBAB(value)
	}
	whitespaceCaseTransformer = func(value string) string {
		return strcase.ToCase(value, strcase.Original, ' ')
	}
)

func getTransformStrategy(c *config.Options) func(string) string {
	switch c.TransformStrategy {
	case "camel":
		return camelCaseTransformer
	case "pascal":
		return pascalCaseTransformer
	case "kebab":
		return kebabCaseTransformer
	case "snake":
		return snakeCaseTransformer
	case "lower":
		return lowerCaseTransformer
	case "upper":
		return upperCaseTransformer
	case "upper-kebab":
		return upperKebabCaseTransformer
	case "upper-snake":
		return upperSnakeCaseTransformer
	case "whitespace":
		return whitespaceCaseTransformer
	default:
		return noopCaseTransformer
	}
}

func (renderUtil) renderSerializers(buf *bytes.Buffer, ts *TypeSpec) error {
	type TplData struct {
		Name              string
		Serializers       []string
		SupportIgnoreCase bool
		SupportUndefined  bool
	}
	data := TplData{
		Name:              ts.Name,
		Serializers:       ts.Config.Serializers,
		SupportIgnoreCase: ts.Config.SupportedFeatures.Contains(config.SupportIgnoreCase),
		SupportUndefined:  ts.Config.SupportedFeatures.Contains(config.SupportUndefined),
	}

	if err := enumTpl.ExecuteTemplate(buf, "enum.serializers.go.tpl", map[string]any{"Type": data}); err != nil {
		return err
	}
	return nil
}

func (renderUtil) renderEntInterfaceSupport(buf *bytes.Buffer, ts *TypeSpec) error {
	type TplData struct {
		Name                string
		SupportEntInterface bool
	}
	data := TplData{
		Name:                ts.Name,
		SupportEntInterface: ts.Config.SupportedFeatures.Contains(config.SupportEntInterface),
	}
	if err := enumTpl.ExecuteTemplate(buf, "enum.misc.ent.go.tpl", map[string]any{"Type": data}); err != nil {
		return err
	}
	return nil
}
