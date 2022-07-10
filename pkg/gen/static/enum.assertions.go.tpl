{{- /* Declare compile time assertion of enum set constants */ -}}
{{- with $ts := .Type -}}
{{- if not $ts.IsFromCsvSource }}
// _{{ $ts.Name }}NoOp is a compile time assertion.
// An "invalid argument/out of bounds" compiler error signifies that the enum values have changed.
// Re-run the enumer command to generate an updated version of {{ $ts.Name }}.
func _{{ $ts.Name }}NoOp() {
	var x [1]struct{}
{{- range $v := $ts.Values }}
	_ = x[{{ $v.ConstName }}-({{ $v.Value }})]
{{- end }}
}

{{ end -}}
{{ end -}}
