{{- /* Declaration of enum's base constants */ -}}
{{- with $ts := .Type -}}
const (
	_{{ $ts.Name }}String      = "{{ $ts.AggregatedValueStrings }}"
	_{{ $ts.Name }}LowerString = "{{ lower $ts.AggregatedValueStrings }}"
{{- if $ts.RequiresGeneratedUndefinedValue }}

	// {{ $ts.Name }}Undefined is the generated zero value of the {{ $ts.Name }} enum.
	{{ $ts.Name }}Undefined {{ $ts.Name }} = 0
{{- end }}
)

{{ end -}}
