//go:generate go run .

package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"text/template"
	"time"

	"github.com/gahojin/go-holiday-japanese/internal"
	"gopkg.in/yaml.v3"
)

type HolidayDetail struct {
	Date   time.Time `yaml:"date"`
	Name   string    `yaml:"name"`    // 祝日名 (日本語)
	NameEn string    `yaml:"name_en"` // 祝日名 (英語)
}

type Holidays map[string]HolidayDetail

type storeMapping struct {
	Diff  uint8
	Index uint8
}

type storeData struct {
	Names   []string
	Mapping []storeMapping
}

var (
	datasetTemplate = template.Must(template.ParseFiles(filepath.Join("dataset.tpl")))

	srcFile = filepath.Join("..", "dataset", "holidays_detailed.yml")
	outFile = filepath.Join("..", "dataset.go")
)

func main() {
	dataset, err := parse(srcFile)
	if err != nil {
		panic(err)
	}

	storeData, err := convert(dataset)
	if err != nil {
		panic(err)
	}

	if err = generate(storeData, outFile); err != nil {
		panic(err)
	}
}

// parse はYAMLファイルを読み込み、祝日情報を返す
func parse(filename string) ([]HolidayDetail, error) {
	file, err := os.ReadFile(filepath.Clean(filename))
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	buf := bytes.NewBuffer(file)

	var dataset Holidays
	if err := yaml.NewDecoder(buf).Decode(&dataset); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	// mapの要素（値）のみを格納するスライス
	holidayList := make([]HolidayDetail, 0, len(dataset))

	// mapをループして、値を取り出しスライスに追加
	for _, detail := range dataset {
		holidayList = append(holidayList, detail)
	}

	sort.Slice(holidayList, func(i, j int) bool {
		a := holidayList[i]
		b := holidayList[j]
		return a.Date.Unix() < b.Date.Unix()
	})

	return holidayList, nil
}

// generate code from template and master data.
func generate(data *storeData, filename string) (err error) {
	buf := new(bytes.Buffer)
	err = datasetTemplate.Execute(buf, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return os.WriteFile(filename, buf.Bytes(), 0600)
}

func convert(dataset []HolidayDetail) (*storeData, error) {
	nameToIndexMap := make(map[string]uint8)
	names := make([]string, 0)
	mapping := make([]storeMapping, 0, len(dataset))

	prevDay := uint32(0)
	for _, info := range dataset {
		date := info.Date
		nameJa := info.Name
		nameEn := info.NameEn

		key := fmt.Sprintf("%s##%s", nameJa, nameEn)
		index, ok := nameToIndexMap[key]
		if !ok {
			n := len(names)
			if n > 254 {
				panic("too many names")
			}
			index = uint8(n)
			names = append(names, nameJa, nameEn)
			nameToIndexMap[key] = index
		}
		epochDay, ok := internal.ToEpochDay(date)
		if !ok {
			return nil, fmt.Errorf("failed to convert date to epoch day: %s", date)
		}
		diff := epochDay - prevDay
		// TODO: diffのoverflowチェックがない
		prevDay = epochDay
		mapping = append(mapping, storeMapping{Diff: uint8(diff), Index: index})
	}

	return &storeData{
		Names:   names,
		Mapping: mapping,
	}, nil
}
