//go:generate go run .

package main

import (
	"bufio"
	"fmt"
	"io"
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

func main() {
	src := filepath.Join("..", "dataset", "holidays_detailed.yml")
	out := filepath.Join("..", "dataset.go")

	dataset, err := parse(src)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(out)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	storeData := convert(dataset)

	if err = generate(storeData, file); err != nil {
		panic(err)
	}
}

// parse はYAMLファイルを読み込み、祝日情報を返す
func parse(filename string) ([]HolidayDetail, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var dataset Holidays
	if err := yaml.NewDecoder(file).Decode(&dataset); err != nil {
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
func generate(data *internal.StoreData, w io.Writer) error {
	writer := bufio.NewWriter(w)
	defer writer.Flush()

	tmpl, err := template.ParseFiles("dataset.tpl")
	if err != nil {
		return err
	}
	return tmpl.Execute(writer, data)
}

func convert(dataset []HolidayDetail) *internal.StoreData {
	nameToIndexMap := make(map[string]uint8)
	names := make([]string, 0)
	mapping := make([]internal.StoreMapping, 0, len(dataset))

	prevDay := uint(0)
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
			panic(fmt.Sprintf("invalid date: %v", date))
		}
		diff := epochDay - prevDay
		prevDay = epochDay
		mapping = append(mapping, internal.StoreMapping{Diff: uint8(diff), Index: index})
	}

	return &internal.StoreData{
		Names:   names,
		Mapping: mapping,
	}
}
