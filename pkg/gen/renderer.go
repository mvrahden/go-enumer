package gen

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/ettle/strcase"
	"github.com/mvrahden/go-enumer/about"
	"github.com/mvrahden/go-enumer/config"
)

type renderer struct {
	cfg  *config.Options
	util renderutil
}

type stringCaseTransformer func(v string) string

func NewRenderer(cfg *config.Options) *renderer {
	r := renderer{}
	for _, fn := range []renderOpt{withDefaults, withTransformStrategy} {
		fn(cfg, &r)
	}
	return &r
}

type renderOpt func(c *config.Options, r *renderer)

func withDefaults(cfg *config.Options, r *renderer) {
	r.cfg = cfg
	r.util = renderutil{
		noopCaseTransformer,
	}
}

func withTransformStrategy(c *config.Options, r *renderer) {
	switch c.TransformStrategy {
	case "camel":
		r.util.transform = camelCaseTransformer
	case "kebab":
		r.util.transform = kebabCaseTransformer
	case "pascal":
		r.util.transform = pascalCaseTransformer
	case "snake":
		r.util.transform = snakeCaseTransformer
	}
}

func (r *renderer) Render(f *File) ([]byte, error) {
	buf := new(bytes.Buffer)
	// write code generation notice and package
	buf.WriteString("// Code generated by \"")
	buf.WriteString(about.ShortInfo())
	buf.WriteString("\"; DO NOT EDIT.\n\n")
	buf.WriteString("package ")
	buf.WriteString(f.Header.Package.Name)
	buf.WriteString("\n\n")

	sort.Slice(f.Imports, func(i, j int) bool {
		return strings.Compare(f.Imports[i].Path, f.Imports[j].Path) < 0
	})

	if len(f.ValueSpecs) == 0 { // nothing to generate further
		return buf.Bytes(), nil
	}

	// write imports list
	buf.WriteString("import (\n")
	for _, v := range f.Imports {
		buf.WriteString("\t")
		if len(v.Name) > 0 {
			buf.WriteString(v.Name)
			buf.WriteString(" ")
		}
		buf.WriteString("\"")
		buf.WriteString(v.Path)
		buf.WriteString("\"\n")
	}
	buf.WriteString(")\n\n")

	// write consts
	tempBuf := new(bytes.Buffer)
	var namestringLen int
	{
		for _, v := range f.ValueSpecs {
			tempBuf.WriteString(v.NameString)
		}
		buf.WriteString("const (\n")
		buf.WriteString(fmt.Sprintf("\t_%sString = \"%s\"\n", r.cfg.TypeAliasName, tempBuf.String()))
		buf.WriteString(fmt.Sprintf("\t_%sLowerString = \"%s\"\n", r.cfg.TypeAliasName, strings.ToLower(tempBuf.String())))
		buf.WriteString(")\n")
		namestringLen = tempBuf.Len()
	}

	// write vars
	{
		tempBuf.Reset()
		for idx, acc := 0, 0; idx < len(f.ValueSpecs); idx++ {
			acc += len(f.ValueSpecs[idx].NameString)
			tempBuf.WriteString(fmt.Sprintf(", %d", acc))
		}
		buf.WriteString("var (\n")
		buf.WriteString(fmt.Sprintf("\t_%sIndices = [%d]uint%d{0%s}\n\n", r.cfg.TypeAliasName, len(f.ValueSpecs)+1, usize(namestringLen), tempBuf.String()))
		tempBuf.Reset()

		for idx := 0; idx < len(f.ValueSpecs); idx++ {
			tempBuf.WriteString(fmt.Sprintf("%s, ", f.ValueSpecs[idx].IdentifierName))
		}
		buf.WriteString(fmt.Sprintf("\t_%sValues = []%s{%s}\n\n", r.cfg.TypeAliasName, r.cfg.TypeAliasName, tempBuf.String()))
		tempBuf.Reset()

		for idx, prev := 0, 0; idx < len(f.ValueSpecs); idx++ {
			l := prev + len(f.ValueSpecs[idx].NameString)
			tempBuf.WriteString(fmt.Sprintf("_%sString[%d:%d], ", r.cfg.TypeAliasName, prev, l))
			prev = l
		}
		buf.WriteString(fmt.Sprintf("\t_%sStrings = []string{%s}\n\n", r.cfg.TypeAliasName, tempBuf.String()))
		buf.WriteString(")\n\n")
	}

	{ // standard functions
		buf.WriteString(fmt.Sprintf(`// %[1]sValues returns all values of the enum.
func %[1]sValues() []%[1]s {
	strs := make([]%[1]s, len(_%[1]sValues))
	copy(strs, _%[1]sValues)
	return _%[1]sValues
}

// %[1]sStrings returns a slice of all String values of the enum.
func %[1]sStrings() []string {
	strs := make([]string, len(_%[1]sStrings))
	copy(strs, _%[1]sStrings)
	return strs
}

// IsValid inspects whether the value is valid enum value.
func (_%[2]s %[1]s) IsValid() bool {
	return _%[2]s >= 0 && _%[2]s < %[1]s(len(_%[1]sIndices)-1)
}

// String returns the string of the enum value.
// If the enum value is invalid.
func (_%[2]s %[1]s) String() string {
	if !_%[2]s.IsValid() {
		return fmt.Sprintf("%[1]s(%%d)", _%[2]s)
	}
	return _%[1]sString[_%[1]sIndices[_%[2]s]:_%[1]sIndices[_%[2]s+1]]
}

`, r.cfg.TypeAliasName, strings.ToLower(string(r.cfg.TypeAliasName[0:1]))))

		buf.WriteString(fmt.Sprintf("var (\n\t_%[1]sStringToValueMap = map[string]%[1]s{\n", r.cfg.TypeAliasName))
		for idx, prev := 0, 0; idx < len(f.ValueSpecs); idx++ {
			l := prev + len(f.ValueSpecs[idx].NameString)
			buf.WriteString(fmt.Sprintf("\t_%[1]sString[%[2]d:%[3]d]: %[4]s,\n", r.cfg.TypeAliasName, prev, l, f.ValueSpecs[idx].IdentifierName))
			buf.WriteString(fmt.Sprintf("\t_%[1]sLowerString[%[2]d:%[3]d]: %[4]s,\n", r.cfg.TypeAliasName, prev, l, f.ValueSpecs[idx].IdentifierName))
			prev = l
		}
		buf.WriteString("}\n)\n\n")
		buf.WriteString(fmt.Sprintf(`func %[1]sFromString(raw string) (%[1]s, bool) {
	v, ok := _%[1]sStringToValueMap[raw]
	if !ok {
		return %[1]s(-1), false
	}
	return v, true
}

`, r.cfg.TypeAliasName))
	}

	{
		r.renderSerializers(buf)
	}

	return buf.Bytes(), nil
}

// usize returns the number of bits of the smallest unsigned integer
// type that will hold n. Used to create the smallest possible slice of
// integers to use as indexes into the concatenated strings.
func usize(n int) int {
	switch {
	case n < 1<<8:
		return 8
	case n < 1<<16:
		return 16
	case n < 1<<32:
		return 32
	default:
		return 64
	}
}

type renderutil struct {
	transform stringCaseTransformer
}

var (
	noopCaseTransformer = func(value string) string {
		return strcase.ToSnake(value)
	}
	snakeCaseTransformer = func(value string) string {
		return strcase.ToSnake(value)
	}
	camelCaseTransformer = func(value string) string {
		return strcase.ToCamel(value)
	}
	pascalCaseTransformer = func(value string) string {
		return strcase.ToPascal(value)
	}
	kebabCaseTransformer = func(value string) string {
		return strcase.ToKebab(value)
	}
)

func (r *renderer) renderSerializers(buf *bytes.Buffer) {
	for _, v := range r.cfg.Serializers {
		switch v {
		case "binary":
			r.renderBinarySerializers(buf)
		case "json":
			r.renderJsonSerializers(buf)
		case "text":
			r.renderTextSerializers(buf)
		case "sql":
			r.renderSqlSerializers(buf)
		case "yaml":
			r.renderYamlSerializers(buf)
		}
	}
}

func (r *renderer) renderBinarySerializers(buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf(`// MarshalBinary implements the encoding.BinaryMarshaler interface for %[1]s.
func (_%[2]s %[1]s) MarshalBinary() ([]byte, error) {
	if !_%[2]s.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %%q as %[1]s", _%[2]s.String())
	}
	return []byte(_%[2]s.String()), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for %[1]s.
func (_%[2]s *%[1]s) UnmarshalBinary(text []byte) error {
	str := string(text)
	if len(str) == 0 {
		return fmt.Errorf("%[1]s cannot be derived from empty string")
	}

	var ok bool
	*_%[2]s, ok = %[1]sFromString(str)
	if !ok {
		return fmt.Errorf("Value %%q does not represent a %[1]s", str)
	}
	return nil
}

`, r.cfg.TypeAliasName, strings.ToLower(string(r.cfg.TypeAliasName[0:1]))))
}

func (r *renderer) renderJsonSerializers(buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf(`// MarshalJSON implements the json.Marshaler interface for %[1]s.
func (_%[2]s %[1]s) MarshalJSON() ([]byte, error) {
	if !_%[2]s.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %%q as %[1]s", _%[2]s.String())
	}
	return json.Marshal(_%[2]s.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for %[1]s.
func (_%[2]s *%[1]s) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("%[1]s should be a string, got %%q", data)
	}
	if len(str) == 0 {
		return fmt.Errorf("%[1]s cannot be derived from empty string")
	}

	var ok bool
	*_%[2]s, ok = %[1]sFromString(str)
	if !ok {
		return fmt.Errorf("Value %%q does not represent a %[1]s", str)
	}
	return nil
}

`, r.cfg.TypeAliasName, strings.ToLower(string(r.cfg.TypeAliasName[0:1]))))
}

func (r *renderer) renderTextSerializers(buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf(`// MarshalText implements the encoding.TextMarshaler interface for %[1]s.
func (_%[2]s %[1]s) MarshalText() ([]byte, error) {
	if !_%[2]s.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %%q as %[1]s", _%[2]s.String())
	}
	return []byte(_%[2]s.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for %[1]s.
func (_%[2]s *%[1]s) UnmarshalText(text []byte) error {
	str := string(text)
	if len(str) == 0 {
		return fmt.Errorf("%[1]s cannot be derived from empty string")
	}

	var ok bool
	*_%[2]s, ok = %[1]sFromString(str)
	if !ok {
		return fmt.Errorf("Value %%q does not represent a %[1]s", str)
	}
	return nil
}

`, r.cfg.TypeAliasName, strings.ToLower(string(r.cfg.TypeAliasName[0:1]))))
}

func (r *renderer) renderSqlSerializers(buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf(`func (_%[2]s %[1]s) Value() (driver.Value, error) {
	if !_%[2]s.IsValid() {
		return nil, fmt.Errorf("Cannot serialize invalid value %%q as %[1]s", _%[2]s.String())
	}
	return _%[2]s.String(), nil
}

func (_%[2]s *%[1]s) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value of %[1]s: %%[1]T(%%[1]v)", value)
	}
	if len(str) == 0 {
		return fmt.Errorf("%[1]s cannot be derived from empty string")
	}

	var ok bool
	*_%[2]s, ok = %[1]sFromString(str)
	if !ok {
		return fmt.Errorf("Value %%q does not represent a %[1]s", str)
	}
	return nil
}

`, r.cfg.TypeAliasName, strings.ToLower(string(r.cfg.TypeAliasName[0:1]))))
}

func (r *renderer) renderYamlSerializers(buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf(`// MarshalYAML implements a YAML Marshaler for %[1]s
func (_%[2]s %[1]s) MarshalYAML() (interface{}, error) {
	if !_%[2]s.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %%q as %[1]s", _%[2]s.String())
	}
	return _%[2]s.String(), nil
}

// UnmarshalYAML implements a YAML Unmarshaler for %[1]s
func (_%[2]s *%[1]s) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}
	if len(str) == 0 {
		return fmt.Errorf("%[1]s cannot be derived from empty string")
	}

	var ok bool
	*_%[2]s, ok = %[1]sFromString(str)
	if !ok {
		return fmt.Errorf("Value %%q does not represent a %[1]s", str)
	}
	return nil
}

`, r.cfg.TypeAliasName, strings.ToLower(string(r.cfg.TypeAliasName[0:1]))))
}