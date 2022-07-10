{{- /* Declaration of enum's base constants */ -}}
{{- with $ts := .Type -}}
const (
	_{{ $ts.Name }}String      = "{{ $ts.AggregatedValueStrings }}"
	_{{ $ts.Name }}LowerString = "{{ lower $ts.AggregatedValueStrings }}"
)

{{ end -}}
