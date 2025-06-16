//go:generate go run .

package main

import (
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"github.com/gahojin/go-holiday-japanese"
	"github.com/gahojin/go-holiday-japanese/internal"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type HolidayDetail struct {
	Date   time.Time `yaml:"date"`
	Name   string    `yaml:"name"`    // 祝日名 (日本語)
	NameEn string    `yaml:"name_en"` // 祝日名 (英語)
}

type Holidays map[string]HolidayDetail

func main() {
	src := filepath.Join("..", "dataset", "holidays_detailed.yml")
	out := filepath.Join("..", "dataset.gob.gz")

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
func generate(data *holiday.StoreData, w io.Writer) error {
	gzipWriter := gzip.NewWriter(w)
	defer gzipWriter.Close()

	encoder := gob.NewEncoder(gzipWriter)
	return encoder.Encode(data)
}

func convert(dataset []HolidayDetail) *holiday.StoreData {
	nameToIndexMap := make(map[string]uint8)
	names := make([]holiday.Name, 0)
	mapping := make([]holiday.Mapping, 0, len(dataset))

	for _, info := range dataset {
		date := info.Date
		nameJa := info.Name
		nameEn := info.NameEn

		key := fmt.Sprintf("%s##%s", nameJa, nameEn)
		index, ok := nameToIndexMap[key]
		if !ok {
			index = uint8(len(names))
			names = append(names, holiday.Name{
				Ja: nameJa,
				En: nameEn,
			})
			nameToIndexMap[key] = index
		}
		epochDay := internal.ToEpochDay(date)
		mapping = append(mapping, holiday.Mapping{
			Day:   epochDay,
			Index: index,
		})
	}

	return &holiday.StoreData{
		Names:   names,
		Mapping: mapping,
	}
}
