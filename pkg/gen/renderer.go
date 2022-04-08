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

type renderer struct{}

func NewRenderer(cfg *config.Options) *renderer {
	r := renderer{}
	return &r
}

func newRenderUtil(cfg *config.Options) *renderUtil {
	r := renderUtil{}
	for _, fn := range []renderUtilOpt{withDefaults, withTransformStrategy, withSupportedFeatures} {
		fn(&r, cfg)
	}
	return &r
}

type renderUtilOpt func(r *renderUtil, c *config.Options)

func withDefaults(r *renderUtil, c *config.Options) {
	r.cfg = c
	r.transform = noopCaseTransformer
}

func withTransformStrategy(r *renderUtil, c *config.Options) {
	r.transform = getTransformStrategy(c)
}

func withSupportedFeatures(r *renderUtil, c *config.Options) {
	r.supportUndefined = c.SupportedFeatures.Contains(config.SupportUndefined)
	r.supportEntInterface = c.SupportedFeatures.Contains(config.SupportEntInterface)
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

	if len(f.TypeSpecs) == 0 { // nothing to generate further
		return buf.Bytes(), nil
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

	for _, ts := range f.TypeSpecs {
		r.renderForTypeSpec(buf, ts)
	}
	return buf.Bytes(), nil
}

func (r *renderer) renderForTypeSpec(buf *bytes.Buffer, ts *TypeSpec) {
	if ts.IsFromCsvSource {
		ts.Config.TransformStrategy = "noop"
	}

	util := newRenderUtil(ts.Config)

	for _, v := range ts.ValueSpecs {
		v.EnumValue = util.transform(v.EnumValue)
	}

	var hasGeneratedUndefinedValue bool
	{ // write consts
		tempBuf := new(bytes.Buffer)
		for _, v := range ts.ValueSpecs {
			tempBuf.WriteString(v.EnumValue)
		}
		buf.WriteString("const (\n")
		buf.WriteString(fmt.Sprintf("\t_%sString = \"%s\"\n", ts.Name, tempBuf))
		buf.WriteString(fmt.Sprintf("\t_%sLowerString = \"%s\"\n", ts.Name, strings.ToLower(tempBuf.String())))
		if ts.HasCanonicalValues {
			tempBuf.Reset()
			for _, v := range ts.ValueSpecs {
				tempBuf.WriteString(v.CanonicalValue)
			}
			buf.WriteString(fmt.Sprintf("\t_%sCanonicalValue = \"%s\"\n", ts.Name, tempBuf))
		}
		buf.WriteString(")\n\n")
		if util.supportUndefined {
			var foundZeroValue bool
			for _, v := range ts.ValueSpecs {
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

`, ts.Name))
			}
		}
	}

	{ // write vars
		tempBuf := new(bytes.Buffer)
		buf.WriteString("var (\n")

		lowerBound := ts.ValueSpecs[0].Value
		if util.supportUndefined && lowerBound > 0 {
			lowerBound = 0
		}
		buf.WriteString(fmt.Sprintf("\t_%[1]sValueRange = [2]%[1]s{%d, %d}\n", ts.Name, lowerBound, ts.ValueSpecs[len(ts.ValueSpecs)-1].Value))

		for idx, prev := 0, uint64(0); idx < len(ts.ValueSpecs); idx++ {
			if idx != 0 && prev == ts.ValueSpecs[idx].Value {
				continue
			}
			prev = ts.ValueSpecs[idx].Value
			tempBuf.WriteString(fmt.Sprintf("%s, ", ts.ValueSpecs[idx].ValueString))
		}
		buf.WriteString(fmt.Sprintf("\t_%sValues = []%s{%s}\n", ts.Name, ts.Name, tempBuf))
		tempBuf.Reset()

		for idx, acc, prev := 0, 0, uint64(0); idx < len(ts.ValueSpecs); idx++ {
			_prev := acc
			acc += len(ts.ValueSpecs[idx].EnumValue)
			if idx != 0 && prev == ts.ValueSpecs[idx].Value {
				continue
			}
			prev = ts.ValueSpecs[idx].Value
			tempBuf.WriteString(fmt.Sprintf("_%sString[%d:%d], ", ts.Name, _prev, acc))
		}
		buf.WriteString(fmt.Sprintf("\t_%sStrings = []string{%s}\n", ts.Name, tempBuf))
		tempBuf.Reset()

		if ts.HasCanonicalValues {
			for idx, acc, prev := 0, 0, uint64(0); idx < len(ts.ValueSpecs); idx++ {
				_prev := acc
				acc += len(ts.ValueSpecs[idx].CanonicalValue)
				if idx != 0 && prev == ts.ValueSpecs[idx].Value {
					continue
				}
				prev = ts.ValueSpecs[idx].Value
				tempBuf.WriteString(fmt.Sprintf("_%sCanonicalValue[%d:%d], ", ts.Name, _prev, acc))
			}
			buf.WriteString(fmt.Sprintf("\t_%sCanonicalValues = []string{%s}\n", ts.Name, tempBuf))
		}
		buf.WriteString(")\n\n")
	}

	if !ts.IsFromCsvSource { // compiletime assertion of numeric sequence
		tempBuf := new(bytes.Buffer)
		if hasGeneratedUndefinedValue {
			tempBuf.WriteString(fmt.Sprintf("\t_ = x[%[1]sUndefined-(0)]\n", ts.Name))
		}
		for _, v := range ts.ValueSpecs {
			tempBuf.WriteString(fmt.Sprintf("\t_ = x[%s-(%s)]\n", v.IdentifierName, v.ValueString))
		}

		buf.WriteString(fmt.Sprintf(`// _%[1]sNoOp is a compile time assertion.
// An "invalid argument/out of bounds" compiler error signifies that the enum values have changed.
// Re-run the enumer command to generate an updated version of %[1]s.
func _%[1]sNoOp() {
	var x [1]struct{}
%[2]s}

	`, ts.Name, tempBuf))
	}

	{ // standard functions
		var offset string
		if ts.ValueSpecs[0].Value > 0 {
			offset = fmt.Sprintf("-%d", ts.ValueSpecs[0].Value)
		}
		undefinedGuard := ""
		if util.supportUndefined && hasGeneratedUndefinedValue {
			undefinedGuard = fmt.Sprintf(`
	if _%[2]s == %[1]sUndefined {
		return ""
	}`, ts.Name, determineReceiverName(ts.Name))
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
	idx := uint(_%[2]s)%[3]s
	return _%[1]sStrings[idx]
}

`, ts.Name, determineReceiverName(ts.Name), offset, undefinedGuard))

		if ts.HasCanonicalValues { // string method for canonical strings
			buf.WriteString(fmt.Sprintf(`// CanonicalValue returns the canonical string of the enum value.
// If the enum value is invalid, it will produce a string
// of the following pattern %[1]s(%%d) instead.
func (_%[2]s %[1]s) CanonicalValue() string {
	if !_%[2]s.IsValid() {
		return fmt.Sprintf("%[1]s(%%d)", _%[2]s)
	}
	idx := uint(_%[2]s)
	return _%[1]sCanonicalValues[idx]
}

`, ts.Name, determineReceiverName(ts.Name)))
		}

		{ // string to numeric value mappings
			buf.WriteString(fmt.Sprintf("var (\n\t_%[1]sStringToValueMap = map[string]%[1]s{\n", ts.Name))
			for idx, prev := 0, 0; idx < len(ts.ValueSpecs); idx++ {
				l := prev + len(ts.ValueSpecs[idx].EnumValue)
				value := ts.ValueSpecs[idx].IdentifierName
				if ts.IsFromCsvSource {
					value = ts.ValueSpecs[idx].ValueString
				}
				buf.WriteString(fmt.Sprintf("\t_%[1]sString[%[2]d:%[3]d]: %[4]s,\n", ts.Name, prev, l, value))
				prev = l
			}
			buf.WriteString("}\n")

			buf.WriteString(fmt.Sprintf("\t_%[1]sLowerStringToValueMap = map[string]%[1]s{\n", ts.Name))
			for idx, prev := 0, 0; idx < len(ts.ValueSpecs); idx++ {
				l := prev + len(ts.ValueSpecs[idx].EnumValue)
				value := ts.ValueSpecs[idx].IdentifierName
				if ts.IsFromCsvSource {
					value = ts.ValueSpecs[idx].ValueString
				}

				buf.WriteString(fmt.Sprintf("\t_%[1]sLowerString[%[2]d:%[3]d]: %[4]s,\n", ts.Name, prev, l, value))
				prev = l
			}
			buf.WriteString("}\n)\n\n")
		}
		zeroValueGuard := ""
		if util.supportUndefined {
			zeroValueGuard = fmt.Sprintf(`
	if len(raw) == 0 {
		return %s(0), true
	}`, ts.Name)
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

`, ts.Name, zeroValueGuard))
	}

	{
		util.renderSerializers(buf, ts)
	}
	{
		if util.supportEntInterface {
			util.renderEntInterfaceSupport(buf, ts)
		}
	}
}

func determineReceiverName(value string) string {
	for _, v := range value {
		return strings.ToLower(string(v))
	}
	return "x"
}

type renderUtil struct {
	cfg                 *config.Options
	transform           stringCaseTransformer
	supportUndefined    bool
	supportEntInterface bool
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

func (r *renderUtil) renderSerializers(buf *bytes.Buffer, ts *TypeSpec) {
	ignoreCase := r.cfg.SupportedFeatures.Contains(config.SupportIgnoreCase)
	for _, v := range r.cfg.Serializers {
		switch v {
		case "binary":
			r.renderBinarySerializers(buf, ts, ignoreCase)
		case "gql":
			r.renderGqlSerializers(buf, ts, ignoreCase)
		case "json":
			r.renderJsonSerializers(buf, ts, ignoreCase)
		case "text":
			r.renderTextSerializers(buf, ts, ignoreCase)
		case "sql":
			r.renderSqlSerializers(buf, ts, ignoreCase)
		case "yaml", "yaml.v3":
			isYamlV3 := v == "yaml.v3"
			r.renderYamlSerializers(buf, ts, ignoreCase, isYamlV3)
		}
	}
}

func (r *renderUtil) renderBinarySerializers(buf *bytes.Buffer, ts *TypeSpec, ignoreCase bool) {
	lookupMethod := "FromString"
	if ignoreCase {
		lookupMethod = "FromStringIgnoreCase"
	}
	zeroValueGuard := ""
	if !r.supportUndefined {
		zeroValueGuard = fmt.Sprintf(`
	if len(str) == 0 {
		return fmt.Errorf("%[1]s cannot be derived from empty string")
	}`, ts.Name)
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

`, ts.Name, determineReceiverName(ts.Name), lookupMethod, zeroValueGuard))
}

func (r *renderUtil) renderGqlSerializers(buf *bytes.Buffer, ts *TypeSpec, ignoreCase bool) {
	lookupMethod := "FromString"
	if ignoreCase {
		lookupMethod = "FromStringIgnoreCase"
	}
	zeroValueGuard := ""
	if !r.supportUndefined {
		zeroValueGuard = fmt.Sprintf(`
	if len(str) == 0 {
		return fmt.Errorf("%[1]s cannot be derived from empty string")
	}`, ts.Name)
	}
	buf.WriteString(fmt.Sprintf(`// MarshalGQL implements the graphql.Marshaler interface for %[1]s.
func (_%[2]s %[1]s) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(_%[2]s.String()))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface for %[1]s.
func (_%[2]s *%[1]s) UnmarshalGQL(value interface{}) error {
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

`, ts.Name, determineReceiverName(ts.Name), lookupMethod, zeroValueGuard))
}

func (r *renderUtil) renderJsonSerializers(buf *bytes.Buffer, ts *TypeSpec, ignoreCase bool) {
	lookupMethod := "FromString"
	if ignoreCase {
		lookupMethod = "FromStringIgnoreCase"
	}
	zeroValueGuard := ""
	if !r.supportUndefined {
		zeroValueGuard = fmt.Sprintf(`
	if len(str) == 0 {
		return fmt.Errorf("%[1]s cannot be derived from empty string")
	}`, ts.Name)
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

`, ts.Name, determineReceiverName(ts.Name), lookupMethod, zeroValueGuard))
}

func (r *renderUtil) renderTextSerializers(buf *bytes.Buffer, ts *TypeSpec, ignoreCase bool) {
	lookupMethod := "FromString"
	if ignoreCase {
		lookupMethod = "FromStringIgnoreCase"
	}
	zeroValueGuard := ""
	if !r.supportUndefined {
		zeroValueGuard = fmt.Sprintf(`
	if len(str) == 0 {
		return fmt.Errorf("%[1]s cannot be derived from empty string")
	}`, ts.Name)
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

`, ts.Name, determineReceiverName(ts.Name), lookupMethod, zeroValueGuard))
}

func (r *renderUtil) renderSqlSerializers(buf *bytes.Buffer, ts *TypeSpec, ignoreCase bool) {
	lookupMethod := "FromString"
	if ignoreCase {
		lookupMethod = "FromStringIgnoreCase"
	}
	zeroValueGuard := ""
	if !r.supportUndefined {
		zeroValueGuard = fmt.Sprintf(`
	if len(str) == 0 {
		return fmt.Errorf("%[1]s cannot be derived from empty string")
	}`, ts.Name)
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

`, ts.Name, determineReceiverName(ts.Name), lookupMethod, zeroValueGuard))
}

func (r *renderUtil) renderYamlSerializers(buf *bytes.Buffer, ts *TypeSpec, ignoreCase, isV3 bool) {
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
`, ts.Name, determineReceiverName(ts.Name)))

	zeroValueGuard := ""
	if !r.supportUndefined {
		zeroValueGuard = fmt.Sprintf(`
	if len(str) == 0 {
		return fmt.Errorf("%[1]s cannot be derived from empty string")
	}`, ts.Name)
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

`, ts.Name, determineReceiverName(ts.Name), lookupMethod, zeroValueGuard))
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

`, ts.Name, determineReceiverName(ts.Name), lookupMethod, zeroValueGuard))
}

func (r *renderUtil) renderEntInterfaceSupport(buf *bytes.Buffer, ts *TypeSpec) {
	buf.WriteString(fmt.Sprintf(`
// Values returns a slice of all String values of the enum.
func (%[1]s) Values() []string {
	return %[1]sStrings()
}
`, ts.Name))
}
