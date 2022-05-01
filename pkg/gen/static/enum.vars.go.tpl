{{/* Declaration of enum's base variables */ -}}
{{- with $ts := .Type -}}
var (
	_{{ $ts.Name }}Values     = [{{ $ts.CountUniqueValues }}]{{ $ts.Name }}{
		{{- range $idx, $v := $ts.Values }}
			{{- if $v.IsAlternativeValue }}{{continue}}{{ end -}}
			{{- $isNotLast := sub (len $ts.Values) 1 | lt $idx -}}
			{{- $v.Value }}
			{{- if $isNotLast }}, {{ end -}}
		{{- end -}}}
	_{{ $ts.Name }}Strings    = [{{ $ts.CountUniqueValues }}]string{
	{{- range $idx, $v := $ts.Values }}
		{{- if $v.IsAlternativeValue }}{{continue}}{{ end -}}
		{{- $isNotLast := sub (len $ts.Values) 1 | lt $idx -}}
		_{{ $ts.Name }}String[{{ $v.Position }}:{{ add $v.Position $v.Length }}]
		{{- if $isNotLast }}, {{ end -}}
	{{- end -}}}
{{- /* Declaration of enum's additional data */ -}}
{{- if $ts.HasAdditionalData }}
	_{{ $ts.Name }}AdditionalData  = [{{ $ts.CountUniqueValues }}]struct{
	{{- range $h := $ts.AdditionalData.Headers }}
		{{ pascal $h.Name }} {{ $h.Type }}
	{{- end }}
	}{
	{{- range $ridx, $r := $ts.AdditionalData.Rows }}
		{
		{{- range $cidx, $c := $r }}
			{{- $isNotLastCell := sub (len $r) 1 | lt $cidx }}
			{{- $c.Value }}{{ if $isNotLastCell }}, {{ end }}
		{{- end -}}
		},
	{{- end }}
	}
{{- end }}
)

{{ end -}}
