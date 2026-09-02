module compatible

go 1.23

replace github.com/gahojin/go-holiday-japanese => ../

require (
	github.com/gahojin/go-holiday-japanese v0.0.0-00010101000000-000000000000
	github.com/holiday-jp/holiday_jp-go v0.0.0-20220125203534-53124b4cc19c
	github.com/stretchr/testify v1.12.1
)

require go.yaml.in/yaml/v3 v3.0.5 // indirect
