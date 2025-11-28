package holiday

var holidayNames = []string{
{{range $name := .Names}}  "{{ $name }}",
{{end}}}
const holidayMapping = "{{range $mapping := .Mapping}}\x{{ printf "%02x" $mapping.Diff }}\x{{ printf "%02x" $mapping.Index }}{{end}}"
