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
	"github.com/mvrahden/go-enumer/pkg/enumer"
	"github.com/mvrahden/go-enumer/pkg/utils/slices"
)

//go:embed static
var templates embed.FS

var (
	headerTpl = template.Must(template.New("header").ParseFS(templates, "static/header.*"))
	enumTpl   = template.Must(template.New("enum").Funcs(tplFuncs).ParseFS(templates, "static/enum.*"))
)

type renderer struct{}

func NewRenderer(cfg *config.Options) *renderer {
	r := renderer{}
	return &r
}

func (r *renderer) Render(f *File) ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := r.renderFileHeader(buf, f); err != nil {
		return nil, fmt.Errorf("failed rendering file header. err: %w", err)
	}

	idx, err := slices.RangeErr(f.TypeSpecs, func(ts *enumer.EnumType, _ int) error {
		return r.renderForTypeSpec(buf, ts)
	})
	if err != nil {
		return nil, fmt.Errorf("failed rendering sources for %q. err: %w", f.TypeSpecs[idx].Name().Name, err)
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
		RepoName:          about.ShortInfo(),
		PackageName:       f.Header.Package.Name,
		Imports:           f.Imports,
		ContainsErrorsPkg: slices.Any(f.Imports, func(v *Import, idx int) bool { return v.Path == "errors" }),
	}
	return headerTpl.ExecuteTemplate(buf, "header.go.tpl", map[string]any{"Header": data})
}

func (r *renderer) renderForTypeSpec(buf *bytes.Buffer, ts *enumer.EnumType) error {
	if ts.HasFileSpec() {
		ts.Config.Options.TransformStrategy = "noop"
	}

	util := newRenderUtil(ts.Config.Options)

	for _, v := range ts.Spec.Values {
		v.EnumValue = util.transform(v.EnumValue)
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
		AdditionalData                  *enumer.AdditionalData
	}
	enum := Enum{
		Name: ts.Name().Name,
		Values: slices.Map(ts.Spec.Values, func(v *enumer.EnumTypeSpecValue, idx int) EnumValue {
			var constName string
			if ts.HasSimpleBlockSpec() {
				constName = v.ConstSpec.Node.Names[0].Name
			}
			return EnumValue{
				Value:     v.ID,
				String:    v.EnumValue,
				ConstName: constName,
				Position: slices.Reduce(ts.Spec.Values[0:idx], func(v *enumer.EnumTypeSpecValue, acc int) int {
					return acc + len(v.EnumValue)
				}),
				Length:             len(v.EnumValue),
				IsAlternativeValue: v.IsAlternative,
			}
		}),
		RequiresGeneratedUndefinedValue: ts.Config.Options.SupportedFeatures.Contains(config.SupportUndefined) &&
			slices.None(ts.Spec.Values, func(v *enumer.EnumTypeSpecValue, idx int) bool { return v.ID == 0 }),
		IsFromCsvSource:   ts.HasFileSpec(),
		HasAdditionalData: ts.Spec.AdditionalData != nil,
		AdditionalData:    ts.Spec.AdditionalData,
	}

	{ // write consts
		type TplData struct {
			Enum
			AggregatedValueStrings string
		}

		data := TplData{
			Enum: enum,
			AggregatedValueStrings: slices.ReduceSeed(ts.Spec.Values, &bytes.Buffer{}, func(v *enumer.EnumTypeSpecValue, acc *bytes.Buffer) *bytes.Buffer {
				acc.WriteString(v.EnumValue)
				return acc
			}).String(),
		}
		if err := enumTpl.ExecuteTemplate(buf, "enum.consts.go.tpl", map[string]any{"Type": data}); err != nil {
			return err
		}
	}

	{ // write vars
		type TplData struct {
			Enum
			CountUniqueValues int // hint: count of all enums, less the alternative values
		}

		data := TplData{
			Enum: enum,
			CountUniqueValues: slices.Count(ts.Spec.Values, func(v *enumer.EnumTypeSpecValue, idx int) bool {
				if idx == 0 {
					return true
				}
				return v.ID != ts.Spec.Values[idx-1].ID
			}),
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
		type Extent struct {
			Min uint64 // hint: the lower numerical bound of the enum set
			Max uint64 // hint: the upper numerical bound of the enum set
		}
		type TplData struct {
			Enum
			Extent         Extent // hint: extent/range of the enum set [min,max]
			RequiresOffset bool
		}

		lowerBound := ts.Spec.Values[0].ID
		if enum.RequiresGeneratedUndefinedValue {
			lowerBound = 0
		}

		data := TplData{
			Enum: enum,
			Extent: Extent{
				Min: lowerBound,
				Max: ts.Spec.Values[len(ts.Spec.Values)-1].ID,
			},
			RequiresOffset: ts.Spec.Values[0].ID > 0,
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
			SupportUndefined: ts.Config.Options.SupportedFeatures.Contains(config.SupportUndefined),
		}

		if err := enumTpl.ExecuteTemplate(buf, "enum.lookup-funcs.go.tpl", map[string]any{"Type": data}); err != nil {
			return err
		}
	}

	{ // serializers
		type TplData struct {
			Name                            string
			RequiresGeneratedUndefinedValue bool
			Serializers                     []string
			SupportIgnoreCase               bool
			SupportUndefined                bool
		}
		data := TplData{
			Name:                            enum.Name,
			RequiresGeneratedUndefinedValue: enum.RequiresGeneratedUndefinedValue,
			Serializers:                     ts.Config.Options.Serializers,
			SupportIgnoreCase:               ts.Config.Options.SupportedFeatures.Contains(config.SupportIgnoreCase),
			SupportUndefined:                ts.Config.Options.SupportedFeatures.Contains(config.SupportUndefined),
		}

		if err := enumTpl.ExecuteTemplate(buf, "enum.serializers.go.tpl", map[string]any{"Type": data}); err != nil {
			return err
		}
	}

	{ // misc (Ent Interface)
		type TplData struct {
			Name                string
			SupportEntInterface bool
		}
		data := TplData{
			Name:                ts.Name().Name,
			SupportEntInterface: ts.Config.Options.SupportedFeatures.Contains(config.SupportEntInterface),
		}
		if err := enumTpl.ExecuteTemplate(buf, "enum.misc.ent.go.tpl", map[string]any{"Type": data}); err != nil {
			return err
		}
	}
	return nil
}

var tplFuncs = template.FuncMap{
	"add": func(a, b int) int {
		return a + b
	},
	"type": enumer.TypeToString,
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
	pascalCaseTransformer     = strcase.ToPascal
	camelCaseTransformer      = strcase.ToCamel
	kebabCaseTransformer      = strcase.ToKebab
	snakeCaseTransformer      = strcase.ToSnake
	lowerCaseTransformer      = strings.ToLower
	upperCaseTransformer      = strings.ToUpper
	upperSnakeCaseTransformer = strcase.ToSNAKE
	upperKebabCaseTransformer = strcase.ToKEBAB
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
