{{- /* Declare ent interface for enum type */ -}}
{{- with $ts := .Type -}}
{{- if $ts.SupportEntInterface -}}
// Values returns a slice of all String values of the enum.
func ({{ $ts.Name }}) Values() []string {
	return {{ $ts.Name }}Strings()
}

{{ end -}}
{{ end -}}
