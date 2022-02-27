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
	for _, fn := range []renderOpt{withDefaults, withTransformStrategy, withSupportedFeatures} {
		fn(cfg, &r)
	}
	return &r
}

type renderOpt func(c *config.Options, r *renderer)

func withDefaults(cfg *config.Options, r *renderer) {
	r.cfg = cfg
	r.util = renderutil{
		transform: noopCaseTransformer,
	}
}

func withTransformStrategy(c *config.Options, r *renderer) {
	switch c.TransformStrategy {
	case "camel":
		r.util.transform = camelCaseTransformer
	case "pascal":
		r.util.transform = pascalCaseTransformer
	case "kebab":
		r.util.transform = kebabCaseTransformer
	case "snake":
		r.util.transform = snakeCaseTransformer
	case "lower":
		r.util.transform = lowerCaseTransformer
	case "upper":
		r.util.transform = upperCaseTransformer
	case "upper-kebab":
		r.util.transform = upperKebabCaseTransformer
	case "upper-snake":
		r.util.transform = upperSnakeCaseTransformer
	case "whitespace":
		r.util.transform = whitespaceCaseTransformer
	}
}

func withSupportedFeatures(c *config.Options, r *renderer) {
	r.util.supportUndefined = c.SupportedFeatures.Contains(config.SupportUndefined)
	r.util.supportEntInterface = c.SupportedFeatures.Contains(config.SupportEntInterface)
}

func (r *renderer) Render(f *File) ([]byte, error) {
	buf := new(bytes.Buffer)
	{ // write code generation notice and package
		buf.WriteString("// Code generated by \"")
		buf.WriteString(about.ShortInfo())
		buf.WriteString("\"; DO NOT EDIT.\n\n")
		buf.WriteString("package ")
		buf.WriteString(f.Header.Package.Name)
		buf.WriteString("\n\n")
	}

	sort.Slice(f.Imports, func(i, j int) bool {
		return strings.Compare(f.Imports[i].Path, f.Imports[j].Path) < 0
	})

	if len(f.ValueSpecs) == 0 { // nothing to generate further
		return buf.Bytes(), nil
	}

	for _, v := range f.ValueSpecs {
		v.EnumString = r.util.transform(v.EnumString)
	}

	{ // write imports list
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
	}

	var hasGeneratedUndefinedValue bool
	{ // write consts
		tempBuf := new(bytes.Buffer)
		for _, v := range f.ValueSpecs {
			tempBuf.WriteString(v.EnumString)
		}
		buf.WriteString("const (\n")
		buf.WriteString(fmt.Sprintf("\t_%sString = \"%s\"\n", r.cfg.TypeAliasName, tempBuf))
		buf.WriteString(fmt.Sprintf("\t_%sLowerString = \"%s\"\n", r.cfg.TypeAliasName, strings.ToLower(tempBuf.String())))
		buf.WriteString(")\n\n")
		if r.util.supportUndefined {
			var foundZeroValue bool
			for _, v := range f.ValueSpecs {
				if v.Value == 0 {
					foundZeroValue = true
					break
				}
			}
			if !foundZeroValue {
				hasGeneratedUndefinedValue = true
				buf.WriteString("")
				buf.WriteString(fmt.Sprintf(`const (
// %[1]sUndefined is the generated zero value of the %[1]s enum.
%[1]sUndefined %[1]s = 0
)

`, r.cfg.TypeAliasName))
			}
		}
	}

	{ // write vars
		tempBuf := new(bytes.Buffer)
		buf.WriteString("var (\n")

		lowerBound := f.ValueSpecs[0].Value
		if r.util.supportUndefined && lowerBound > 0 {
			lowerBound = 0
		}
		buf.WriteString(fmt.Sprintf("\t_%[1]sValueRange = [2]%[1]s{%d, %d}\n", r.cfg.TypeAliasName, lowerBound, f.ValueSpecs[len(f.ValueSpecs)-1].Value))

		for idx, prev := 0, uint64(0); idx < len(f.ValueSpecs); idx++ {
			if idx != 0 && prev == f.ValueSpecs[idx].Value {
				continue
			}
			prev = f.ValueSpecs[idx].Value
			tempBuf.WriteString(fmt.Sprintf("%s, ", f.ValueSpecs[idx].ValueString))
		}
		buf.WriteString(fmt.Sprintf("\t_%sValues = []%s{%s}\n", r.cfg.TypeAliasName, r.cfg.TypeAliasName, tempBuf))
		tempBuf.Reset()

		for idx, acc, prev := 0, 0, uint64(0); idx < len(f.ValueSpecs); idx++ {
			_prev := acc
			acc += len(f.ValueSpecs[idx].EnumString)
			if idx != 0 && prev == f.ValueSpecs[idx].Value {
				continue
			}
			prev = f.ValueSpecs[idx].Value
			tempBuf.WriteString(fmt.Sprintf("_%sString[%d:%d], ", r.cfg.TypeAliasName, _prev, acc))
		}
		buf.WriteString(fmt.Sprintf("\t_%sStrings = []string{%s}\n", r.cfg.TypeAliasName, tempBuf))
		buf.WriteString(")\n\n")
	}

	{ // compiletime assertion of numeric sequence
		tempBuf := new(bytes.Buffer)
		if hasGeneratedUndefinedValue {
			tempBuf.WriteString(fmt.Sprintf("\t_ = x[%[1]sUndefined-(0)]\n", r.cfg.TypeAliasName))
		}
		for _, v := range f.ValueSpecs {
			tempBuf.WriteString(fmt.Sprintf("\t_ = x[%s-(%s)]\n", v.IdentifierName, v.ValueString))
		}

		buf.WriteString(fmt.Sprintf(`// _%[1]sNoOp is a compile time assertion.
// An "invalid argument/out of bounds" compiler error signifies that the enum values have changed.
// Re-run the enumer command to generate an updated version of %[1]s.
func _%[1]sNoOp() {
	var x [1]struct{}
%[2]s}

	`, r.cfg.TypeAliasName, tempBuf))
	}

	{ // standard functions
		var offset string
		if f.ValueSpecs[0].Value > 0 {
			offset = fmt.Sprintf("-%d", f.ValueSpecs[0].Value)
		}
		undefinedGuard := ""
		if r.util.supportUndefined && hasGeneratedUndefinedValue {
			undefinedGuard = fmt.Sprintf(`
	if _%[2]s == %[1]sUndefined {
		return ""
	}`, r.cfg.TypeAliasName, determineReceiverName(r.cfg.TypeAliasName))
		}
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
	return _%[2]s >= _%[1]sValueRange[0] && _%[2]s <= _%[1]sValueRange[1]
}

// String returns the string of the enum value.
// If the enum value is invalid, it will produce a string
// of the following pattern %[1]s(%%d) instead.
func (_%[2]s %[1]s) String() string {
	if !_%[2]s.IsValid() {
		return fmt.Sprintf("%[1]s(%%d)", _%[2]s)
	}%[4]s
	idx := int(_%[2]s)%[3]s
	return _%[1]sStrings[idx]
}

`, r.cfg.TypeAliasName, determineReceiverName(r.cfg.TypeAliasName), offset, undefinedGuard))

		buf.WriteString(fmt.Sprintf("var (\n\t_%[1]sStringToValueMap = map[string]%[1]s{\n", r.cfg.TypeAliasName))
		for idx, prev := 0, 0; idx < len(f.ValueSpecs); idx++ {
			l := prev + len(f.ValueSpecs[idx].EnumString)
			buf.WriteString(fmt.Sprintf("\t_%[1]sString[%[2]d:%[3]d]: %[4]s,\n", r.cfg.TypeAliasName, prev, l, f.ValueSpecs[idx].IdentifierName))
			prev = l
		}
		buf.WriteString("}\n")

		buf.WriteString(fmt.Sprintf("\t_%[1]sLowerStringToValueMap = map[string]%[1]s{\n", r.cfg.TypeAliasName))
		for idx, prev := 0, 0; idx < len(f.ValueSpecs); idx++ {
			l := prev + len(f.ValueSpecs[idx].EnumString)
			buf.WriteString(fmt.Sprintf("\t_%[1]sLowerString[%[2]d:%[3]d]: %[4]s,\n", r.cfg.TypeAliasName, prev, l, f.ValueSpecs[idx].IdentifierName))
			prev = l
		}
		buf.WriteString("}\n)\n\n")
		zeroValueGuard := ""
		if r.util.supportUndefined {
			zeroValueGuard = fmt.Sprintf(`
	if len(raw) == 0 {
		return %s(0), true
	}`, r.cfg.TypeAliasName)
		}
		buf.WriteString(fmt.Sprintf(`// %[1]sFromString determines the enum value with an exact case match.
func %[1]sFromString(raw string) (%[1]s, bool) {%[2]s
	v, ok := _%[1]sStringToValueMap[raw]
	if !ok {
		return %[1]s(0), false
	}
	return v, true
}

// %[1]sFromStringIgnoreCase determines the enum value with a case-insensitive match.
func %[1]sFromStringIgnoreCase(raw string) (%[1]s, bool) {
	v, ok := %[1]sFromString(raw)
	if ok {
		return v, ok
	}
	v, ok = _%[1]sLowerStringToValueMap[raw]
	if !ok {
		return %[1]s(0), false
	}
	return v, true
}

`, r.cfg.TypeAliasName, zeroValueGuard))
	}

	{
		r.renderSerializers(buf)
	}
	{
		if r.util.supportEntInterface {
			r.renderEntInterfaceSupport(buf)
		}
	}

	return buf.Bytes(), nil
}

func determineReceiverName(value string) string {
	for _, v := range value {
		return strings.ToLower(string(v))
	}
	return "x"
}

type renderutil struct {
	transform           stringCaseTransformer
	supportUndefined    bool
	supportEntInterface bool
}

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

func (r *renderer) renderSerializers(buf *bytes.Buffer) {
	ignoreCase := r.cfg.SupportedFeatures.Contains("ignore-case")
	for _, v := range r.cfg.Serializers {
		switch v {
		case "binary":
			r.renderBinarySerializers(buf, ignoreCase)
		case "gql":
			r.renderGqlSerializers(buf, ignoreCase)
		case "json":
			r.renderJsonSerializers(buf, ignoreCase)
		case "text":
			r.renderTextSerializers(buf, ignoreCase)
		case "sql":
			r.renderSqlSerializers(buf, ignoreCase)
		case "yaml", "yaml.v3":
			isYamlV3 := v == "yaml.v3"
			r.renderYamlSerializers(buf, ignoreCase, isYamlV3)
		}
	}
}

func (r *renderer) renderBinarySerializers(buf *bytes.Buffer, ignoreCase bool) {
	lookupMethod := "FromString"
	if ignoreCase {
		lookupMethod = "FromStringIgnoreCase"
	}
	zeroValueGuard := ""
	if !r.util.supportUndefined {
		zeroValueGuard = fmt.Sprintf(`
	if len(str) == 0 {
		return fmt.Errorf("%[1]s cannot be derived from empty string")
	}`, r.cfg.TypeAliasName)
	}
	buf.WriteString(fmt.Sprintf(`// MarshalBinary implements the encoding.BinaryMarshaler interface for %[1]s.
func (_%[2]s %[1]s) MarshalBinary() ([]byte, error) {
	if !_%[2]s.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %%q as %[1]s", _%[2]s)
	}
	return []byte(_%[2]s.String()), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for %[1]s.
func (_%[2]s *%[1]s) UnmarshalBinary(text []byte) error {
	str := string(text)%[4]s

	var ok bool
	*_%[2]s, ok = %[1]s%[3]s(str)
	if !ok {
		return fmt.Errorf("Value %%q does not represent a %[1]s", str)
	}
	return nil
}

`, r.cfg.TypeAliasName, determineReceiverName(r.cfg.TypeAliasName), lookupMethod, zeroValueGuard))
}

func (r *renderer) renderGqlSerializers(buf *bytes.Buffer, ignoreCase bool) {
	lookupMethod := "FromString"
	if ignoreCase {
		lookupMethod = "FromStringIgnoreCase"
	}
	zeroValueGuard := ""
	if !r.util.supportUndefined {
		zeroValueGuard = fmt.Sprintf(`
	if len(str) == 0 {
		return fmt.Errorf("%[1]s cannot be derived from empty string")
	}`, r.cfg.TypeAliasName)
	}
	buf.WriteString(fmt.Sprintf(`// MarshalGQL implements the graphql.Marshaler interface for %[1]s.
func (_%[2]s %[1]s) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(_%[2]s.String()))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface for %[1]s.
func (_%[2]s *%[1]s) UnmarshalGQL(value interface{}) error {
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
	}%[4]s

	var ok bool
	*_%[2]s, ok = %[1]s%[3]s(str)
	if !ok {
		return fmt.Errorf("Value %%q does not represent a %[1]s", str)
	}
	return nil
}

`, r.cfg.TypeAliasName, determineReceiverName(r.cfg.TypeAliasName), lookupMethod, zeroValueGuard))
}

func (r *renderer) renderJsonSerializers(buf *bytes.Buffer, ignoreCase bool) {
	lookupMethod := "FromString"
	if ignoreCase {
		lookupMethod = "FromStringIgnoreCase"
	}
	zeroValueGuard := ""
	if !r.util.supportUndefined {
		zeroValueGuard = fmt.Sprintf(`
	if len(str) == 0 {
		return fmt.Errorf("%[1]s cannot be derived from empty string")
	}`, r.cfg.TypeAliasName)
	}
	buf.WriteString(fmt.Sprintf(`// MarshalJSON implements the json.Marshaler interface for %[1]s.
func (_%[2]s %[1]s) MarshalJSON() ([]byte, error) {
	if !_%[2]s.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %%q as %[1]s", _%[2]s)
	}
	return json.Marshal(_%[2]s.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for %[1]s.
func (_%[2]s *%[1]s) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("%[1]s should be a string, got %%q", data)
	}%[4]s

	var ok bool
	*_%[2]s, ok = %[1]s%[3]s(str)
	if !ok {
		return fmt.Errorf("Value %%q does not represent a %[1]s", str)
	}
	return nil
}

`, r.cfg.TypeAliasName, determineReceiverName(r.cfg.TypeAliasName), lookupMethod, zeroValueGuard))
}

func (r *renderer) renderTextSerializers(buf *bytes.Buffer, ignoreCase bool) {
	lookupMethod := "FromString"
	if ignoreCase {
		lookupMethod = "FromStringIgnoreCase"
	}
	zeroValueGuard := ""
	if !r.util.supportUndefined {
		zeroValueGuard = fmt.Sprintf(`
	if len(str) == 0 {
		return fmt.Errorf("%[1]s cannot be derived from empty string")
	}`, r.cfg.TypeAliasName)
	}
	buf.WriteString(fmt.Sprintf(`// MarshalText implements the encoding.TextMarshaler interface for %[1]s.
func (_%[2]s %[1]s) MarshalText() ([]byte, error) {
	if !_%[2]s.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %%q as %[1]s", _%[2]s)
	}
	return []byte(_%[2]s.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for %[1]s.
func (_%[2]s *%[1]s) UnmarshalText(text []byte) error {
	str := string(text)%[4]s

	var ok bool
	*_%[2]s, ok = %[1]s%[3]s(str)
	if !ok {
		return fmt.Errorf("Value %%q does not represent a %[1]s", str)
	}
	return nil
}

`, r.cfg.TypeAliasName, determineReceiverName(r.cfg.TypeAliasName), lookupMethod, zeroValueGuard))
}

func (r *renderer) renderSqlSerializers(buf *bytes.Buffer, ignoreCase bool) {
	lookupMethod := "FromString"
	if ignoreCase {
		lookupMethod = "FromStringIgnoreCase"
	}
	zeroValueGuard := ""
	if !r.util.supportUndefined {
		zeroValueGuard = fmt.Sprintf(`
	if len(str) == 0 {
		return fmt.Errorf("%[1]s cannot be derived from empty string")
	}`, r.cfg.TypeAliasName)
	}
	buf.WriteString(fmt.Sprintf(`func (_%[2]s %[1]s) Value() (driver.Value, error) {
	if !_%[2]s.IsValid() {
		return nil, fmt.Errorf("Cannot serialize invalid value %%q as %[1]s", _%[2]s)
	}
	return _%[2]s.String(), nil
}

func (_%[2]s *%[1]s) Scan(value interface{}) error {
	var str string
	switch v := value.(type) {
	case nil:
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value of %[1]s: %%[1]T(%%[1]v)", value)
	}%[4]s

	var ok bool
	*_%[2]s, ok = %[1]s%[3]s(str)
	if !ok {
		return fmt.Errorf("Value %%q does not represent a %[1]s", str)
	}
	return nil
}

`, r.cfg.TypeAliasName, determineReceiverName(r.cfg.TypeAliasName), lookupMethod, zeroValueGuard))
}

func (r *renderer) renderYamlSerializers(buf *bytes.Buffer, ignoreCase, isV3 bool) {
	lookupMethod := "FromString"
	if ignoreCase {
		lookupMethod = "FromStringIgnoreCase"
	}
	buf.WriteString(fmt.Sprintf(`// MarshalYAML implements a YAML Marshaler for %[1]s.
func (_%[2]s %[1]s) MarshalYAML() (interface{}, error) {
	if !_%[2]s.IsValid() {
		return nil, fmt.Errorf("Cannot marshal invalid value %%q as %[1]s", _%[2]s)
	}
	return _%[2]s.String(), nil
}
`, r.cfg.TypeAliasName, determineReceiverName(r.cfg.TypeAliasName)))

	zeroValueGuard := ""
	if !r.util.supportUndefined {
		zeroValueGuard = fmt.Sprintf(`
	if len(str) == 0 {
		return fmt.Errorf("%[1]s cannot be derived from empty string")
	}`, r.cfg.TypeAliasName)
	}
	if isV3 {
		buf.WriteString(fmt.Sprintf(`
// UnmarshalYAML implements a YAML Unmarshaler for %[1]s.
func (_%[2]s *%[1]s) UnmarshalYAML(n *yaml.Node) error {
	const stringTag = "!!str"
	if n.ShortTag() != stringTag {
		return fmt.Errorf("%[1]s must be derived from a string node")
	}
	str := n.Value%[4]s

	var ok bool
	*_%[2]s, ok = %[1]s%[3]s(str)
	if !ok {
		return fmt.Errorf("Value %%q does not represent a %[1]s", str)
	}
	return nil
}

`, r.cfg.TypeAliasName, determineReceiverName(r.cfg.TypeAliasName), lookupMethod, zeroValueGuard))
		return
	}
	buf.WriteString(fmt.Sprintf(`
// UnmarshalYAML implements a YAML Unmarshaler for %[1]s.
func (_%[2]s *%[1]s) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}%[4]s

	var ok bool
	*_%[2]s, ok = %[1]s%[3]s(str)
	if !ok {
		return fmt.Errorf("Value %%q does not represent a %[1]s", str)
	}
	return nil
}

`, r.cfg.TypeAliasName, determineReceiverName(r.cfg.TypeAliasName), lookupMethod, zeroValueGuard))
}

func (r *renderer) renderEntInterfaceSupport(buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf(`
// Values returns a slice of all String values of the enum.
func (%[1]s) Values() []string {
	return %[1]sStrings()
}
`, r.cfg.TypeAliasName))
}
