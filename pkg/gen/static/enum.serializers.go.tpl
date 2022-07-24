{{- /* Declare binary interface for enum type */ -}}
{{- with $ts := .Type -}}
{{- if contains $ts.Serializers "binary" }}
// MarshalBinary implements the encoding.BinaryMarshaler interface for {{ $ts.Name }}.
func ({{ receiver $ts.Name }} {{ $ts.Name }}) MarshalBinary() ([]byte, error) {
	if err := {{ receiver $ts.Name }}.Validate(); err != nil {
		return nil, fmt.Errorf("Cannot marshal value %q as {{ $ts.Name }}. %w", {{ receiver $ts.Name }}, err)
	}
	return []byte({{ receiver $ts.Name }}.String()), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface for {{ $ts.Name }}.
func ({{ receiver $ts.Name }} *{{ $ts.Name }}) UnmarshalBinary(text []byte) error {
	str := string(text)
{{- if not $ts.SupportUndefined }}
	if len(str) == 0 {
		return fmt.Errorf("{{ $ts.Name }} cannot be derived from empty string")
	}
{{- end }}

	var ok bool
	*{{ receiver $ts.Name }}, ok = {{ $ts.Name }}FromString{{ if $ts.SupportIgnoreCase }}IgnoreCase{{ end }}(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a {{ $ts.Name }}", str)
	}
	return nil
}
{{ end }}
{{- if contains $ts.Serializers "bson" }}
// MarshalBSONValue implements the bson.ValueMarshaler interface for {{ $ts.Name }}.
func ({{ receiver $ts.Name }} {{ $ts.Name }}) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if err := {{ receiver $ts.Name }}.Validate(); err != nil {
		return 0, nil, fmt.Errorf("Cannot marshal value %q as {{ $ts.Name }}. %w", {{ receiver $ts.Name }}, err)
	}
{{- if $ts.RequiresGeneratedUndefinedValue }}
	if {{ receiver $ts.Name }} == 0 {
		return bsontype.Undefined, nil, nil
	}
{{- end }}
	return bson.MarshalValue({{ receiver $ts.Name }}.String())
}

// UnmarshalBSONValue implements the bson.ValueUnmarshaler interface for {{ $ts.Name }}.
func ({{ receiver $ts.Name }} *{{ $ts.Name }}) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t != bsontype.String {{- if $ts.SupportUndefined }} && t != bsontype.Undefined {{- end }} {
		return fmt.Errorf("{{ $ts.Name }} should be a string, got %q of Type %q", data, t)
	}
	str, data, ok := bsoncore.ReadString(data)
	if !ok {
		return fmt.Errorf("failed reading value as string, got %q", data)
	}
{{- if not $ts.SupportUndefined }}
	if len(str) == 0 {
		return fmt.Errorf("{{ $ts.Name }} cannot be derived from empty string")
	}
{{- end }}

	*{{ receiver $ts.Name }}, ok = {{ $ts.Name }}FromString{{ if $ts.SupportIgnoreCase }}IgnoreCase{{ end }}(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a {{ $ts.Name }}", str)
	}
	return nil
}
{{ end }}
{{- if contains $ts.Serializers "graphql" }}
// MarshalGQL implements the graphql.Marshaler interface for {{ $ts.Name }}.
func ({{ receiver $ts.Name }} {{ $ts.Name }}) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote({{ receiver $ts.Name }}.String()))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface for {{ $ts.Name }}.
func ({{ receiver $ts.Name }} *{{ $ts.Name }}) UnmarshalGQL(value interface{}) error {
	var str string
	switch v := value.(type) {
	{{- if $ts.SupportUndefined }}
	case nil:
	{{- end }}
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value of {{ $ts.Name }}: %[1]T(%[1]v)", value)
	}
{{- if not $ts.SupportUndefined }}
	if len(str) == 0 {
		return fmt.Errorf("{{ $ts.Name }} cannot be derived from empty string")
	}
{{- end }}

	var ok bool
	*{{ receiver $ts.Name }}, ok = {{ $ts.Name }}FromString{{ if $ts.SupportIgnoreCase }}IgnoreCase{{ end }}(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a {{ $ts.Name }}", str)
	}
	return nil
}
{{ end }}
{{- if contains $ts.Serializers "json" }}
// MarshalJSON implements the json.Marshaler interface for {{ $ts.Name }}.
func ({{ receiver $ts.Name }} {{ $ts.Name }}) MarshalJSON() ([]byte, error) {
	if err := {{ receiver $ts.Name }}.Validate(); err != nil {
		return nil, fmt.Errorf("Cannot marshal value %q as {{ $ts.Name }}. %w", {{ receiver $ts.Name }}, err)
	}
	return json.Marshal({{ receiver $ts.Name }}.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for {{ $ts.Name }}.
func ({{ receiver $ts.Name }} *{{ $ts.Name }}) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("{{ $ts.Name }} should be a string, got %q", data)
	}
{{- if not $ts.SupportUndefined }}
	if len(str) == 0 {
		return fmt.Errorf("{{ $ts.Name }} cannot be derived from empty string")
	}
{{- end }}

	var ok bool
	*{{ receiver $ts.Name }}, ok = {{ $ts.Name }}FromString{{ if $ts.SupportIgnoreCase }}IgnoreCase{{ end }}(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a {{ $ts.Name }}", str)
	}
	return nil
}
{{ end }}
{{- if contains $ts.Serializers "sql" }}
// Value implements the sql/driver.Valuer interface for {{ $ts.Name }}.
func ({{ receiver $ts.Name }} {{ $ts.Name }}) Value() (driver.Value, error) {
	if err := {{ receiver $ts.Name }}.Validate(); err != nil {
		return nil, fmt.Errorf("Cannot serialize value %q as {{ $ts.Name }}. %w", {{ receiver $ts.Name }}, err)
	}
{{- if $ts.RequiresGeneratedUndefinedValue }}
	if {{ receiver $ts.Name }} == 0 {
		return nil, nil
	}
{{- end }}
	return {{ receiver $ts.Name }}.String(), nil
}

// Scan implements the sql/driver.Scanner interface for {{ $ts.Name }}.
func ({{ receiver $ts.Name }} *{{ $ts.Name }}) Scan(value interface{}) error {
	var str string
	switch v := value.(type) {
	{{- if $ts.SupportUndefined }}
	case nil:
	{{- end }}
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value of {{ $ts.Name }}: %[1]T(%[1]v)", value)
	}
{{- if not $ts.SupportUndefined }}
	if len(str) == 0 {
		return fmt.Errorf("{{ $ts.Name }} cannot be derived from empty string")
	}
{{- end }}

	var ok bool
	*{{ receiver $ts.Name }}, ok = {{ $ts.Name }}FromString{{ if $ts.SupportIgnoreCase }}IgnoreCase{{ end }}(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a {{ $ts.Name }}", str)
	}
	return nil
}
{{ end }}
{{- if contains $ts.Serializers "text" }}
// MarshalText implements the encoding.TextMarshaler interface for {{ $ts.Name }}.
func ({{ receiver $ts.Name }} {{ $ts.Name }}) MarshalText() ([]byte, error) {
	if err := {{ receiver $ts.Name }}.Validate(); err != nil {
		return nil, fmt.Errorf("Cannot marshal value %q as {{ $ts.Name }}. %w", {{ receiver $ts.Name }}, err)
	}
	return []byte({{ receiver $ts.Name }}.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for {{ $ts.Name }}.
func ({{ receiver $ts.Name }} *{{ $ts.Name }}) UnmarshalText(text []byte) error {
	str := string(text)
{{- if not $ts.SupportUndefined }}
	if len(str) == 0 {
		return fmt.Errorf("{{ $ts.Name }} cannot be derived from empty string")
	}
{{- end }}

	var ok bool
	*{{ receiver $ts.Name }}, ok = {{ $ts.Name }}FromString{{ if $ts.SupportIgnoreCase }}IgnoreCase{{ end }}(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a {{ $ts.Name }}", str)
	}
	return nil
}
{{ end }}
{{- $serializeYamlV3 := contains $ts.Serializers "yaml.v3" -}}
{{- if or (contains $ts.Serializers "yaml") $serializeYamlV3 }}
// MarshalYAML implements a YAML Marshaler for {{ $ts.Name }}.
func ({{ receiver $ts.Name }} {{ $ts.Name }}) MarshalYAML() (interface{}, error) {
	if err := {{ receiver $ts.Name }}.Validate(); err != nil {
		return nil, fmt.Errorf("Cannot marshal value %q as {{ $ts.Name }}. %w", {{ receiver $ts.Name }}, err)
	}
	return {{ receiver $ts.Name }}.String(), nil
}

{{ if $serializeYamlV3 -}}
// UnmarshalYAML implements a YAML Unmarshaler for {{ $ts.Name }}.
func ({{ receiver $ts.Name }} *{{ $ts.Name }}) UnmarshalYAML(n *yaml.Node) error {
	const stringTag = "!!str"
	if n.ShortTag() != stringTag {
		return fmt.Errorf("{{ $ts.Name }} must be derived from a string node")
	}
	str := n.Value
{{- else -}}
// UnmarshalYAML implements a YAML Unmarshaler for {{ $ts.Name }}.
func ({{ receiver $ts.Name }} *{{ $ts.Name }}) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}
{{- end }}
{{- if not $ts.SupportUndefined }}
	if len(str) == 0 {
		return fmt.Errorf("{{ $ts.Name }} cannot be derived from empty string")
	}
{{- end }}

	var ok bool
	*{{ receiver $ts.Name }}, ok = {{ $ts.Name }}FromString{{ if $ts.SupportIgnoreCase }}IgnoreCase{{ end }}(str)
	if !ok {
		return fmt.Errorf("Value %q does not represent a {{ $ts.Name }}", str)
	}
	return nil
}
{{ end }}

{{ end -}}
