
{{- /* Declare lookup functions of enum type */ -}}
{{- with $ts := .Type -}}
var (
	_{{ $ts.Name }}StringToValueMap = map[string]{{ $ts.Name }}{
{{- range $v := $ts.Values }}
		_{{ $ts.Name }}String[{{ $v.Position }}:{{ add $v.Position $v.Length }}]: {{ if $ts.IsFromCsvSource }}{{ $v.Value }}{{ else }}{{ $v.ConstName }}{{ end }},
{{- end }}
	}
	_{{ $ts.Name }}LowerStringToValueMap = map[string]{{ $ts.Name }}{
{{- range $v := $ts.Values }}
		_{{ $ts.Name }}LowerString[{{ $v.Position }}:{{ add $v.Position $v.Length }}]: {{ if $ts.IsFromCsvSource }}{{ $v.Value }}{{ else }}{{ $v.ConstName }}{{ end }},
{{- end }}
	}
)

// {{ $ts.Name }}FromString determines the enum value with an exact case match.
func {{ $ts.Name }}FromString(raw string) ({{ $ts.Name }}, bool) {
{{- if $ts.SupportUndefined }}
	if len(raw) == 0 {
		return {{ $ts.Name }}(0), true
	}
{{- end }}
	v, ok := _{{ $ts.Name }}StringToValueMap[raw]
	if !ok {
		return {{ $ts.Name }}(0), false
	}
	return v, true
}

// {{ $ts.Name }}FromStringIgnoreCase determines the enum value with a case-insensitive match.
func {{ $ts.Name }}FromStringIgnoreCase(raw string) ({{ $ts.Name }}, bool) {
{{- if $ts.SupportUndefined }}
	if len(raw) == 0 {
		return {{ $ts.Name }}(0), true
	}
{{- end }}
	v, ok := {{ $ts.Name }}FromString(raw)
	if ok {
		return v, ok
	}
	v, ok = _{{ $ts.Name }}LowerStringToValueMap[raw]
	if !ok {
		return {{ $ts.Name }}(0), false
	}
	return v, true
}

{{ end -}}
