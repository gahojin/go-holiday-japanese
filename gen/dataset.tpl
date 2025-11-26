package holiday

import "github.com/gahojin/go-holiday-japanese/internal"

var holidayNames = []string{
{{range $name := .Names}}  "{{ $name }}",
{{end}}}
var holidayMapping = []internal.StoreMapping{
{{range $mapping := .Mapping}}  {Diff: {{$mapping.Diff}}, Index: {{$mapping.Index}}},
{{end}}}
