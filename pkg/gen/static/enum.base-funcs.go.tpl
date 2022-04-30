{{- /* Declare base functions of enum type */ -}}
{{- with $ts := .Type -}}
// {{ $ts.Name }}Values returns all values of the enum.
func {{ $ts.Name }}Values() []{{ $ts.Name }} {
	cp := _{{ $ts.Name }}Values
	return cp[:]
}

// {{ $ts.Name }}Strings returns a slice of all String values of the enum.
func {{ $ts.Name }}Strings() []string {
	cp := _{{ $ts.Name }}Strings
	return cp[:]
}

// IsValid inspects whether the value is valid enum value.
func ({{ receiver $ts.Name }} {{ $ts.Name }}) IsValid() bool {
	return {{ receiver $ts.Name }} >= _{{ $ts.Name }}ValueRange[0] && {{ receiver $ts.Name }} <= _{{ $ts.Name }}ValueRange[1]
}

// Validate whether the value is within the range of enum values.
func ({{ receiver $ts.Name }} {{ $ts.Name }}) Validate() error {
	if !{{ receiver $ts.Name }}.IsValid() {
		return fmt.Errorf("{{ $ts.Name }}(%d) is %w", {{ receiver $ts.Name }}, ErrNoValidEnum)
	}
	return nil
}

// String returns the string of the enum value.
// If the enum value is invalid, it will produce a string
// of the following pattern {{ $ts.Name }}(%d) instead.
func ({{ receiver $ts.Name }} {{ $ts.Name }}) String() string {
	if !{{ receiver $ts.Name }}.IsValid() {
		return fmt.Sprintf("{{ $ts.Name }}(%d)", {{ receiver $ts.Name }})
	}
{{- if $ts.RequiresGeneratedUndefinedValue }}
{{- /* This block assures the proper serialization of the generated undefined Value */}}
	if {{ receiver $ts.Name }} == {{ $ts.Name }}Undefined {
		return ""
	}
{{- end }}
	idx := uint({{ receiver $ts.Name }}){{- if $ts.RequiresOffset }} - 1{{- end }}
	return _{{ $ts.Name }}Strings[idx]
}

{{ if $ts.HasAdditionalData }}
{{- /* Generate typed getter for additional data */}}
{{- range $col := $ts.AdditionalData.Headers -}}
// Get{{ pascal $col.Name }} returns the "{{ $col.Name }}" of the enum value.
func ({{ receiver $ts.Name }} {{ $ts.Name }}) Get{{ pascal $col.Name }}() {{ $col.Type }} {
	idx := uint({{ receiver $ts.Name }}){{- if $ts.RequiresOffset }} - 1{{- end }}
	d := _{{ $ts.Name }}AdditionalData[idx]
	return d.{{ pascal $col.Name }}
}

{{ end -}}
{{- end -}}
{{- end -}}
