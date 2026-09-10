package internal

var HolidayNames = []string{
{{ range $name := .Names }}	"{{ $name }}",
{{ end }}}

const HolidayMapping = "{{ range $mapping := .Mapping}}\x{{ printf "%02x" $mapping.Diff }}\x{{ printf "%02x" $mapping.Index }}{{end}}"

const EpochDayMax = {{ .EpochDayMax }}
