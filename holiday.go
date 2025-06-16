package holiday

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"encoding/gob"
	"fmt"
	"github.com/gahojin/holiday-japanese-go/internal"
	"time"
)

var (
	//go:embed dataset.gob.gz
	rawDataset []byte

	dataset = decodeDatasetCached()
)

func decodeDataset() (*StoreData, error) {
	gzipReader, err := gzip.NewReader(bytes.NewReader(rawDataset))
	if err != nil {
		return nil, fmt.Errorf("failed to decompress dataset: %w", err)
	}

	var data StoreData
	decoder := gob.NewDecoder(gzipReader)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode dataset: %w", err)
	}

	namesLen := len(data.Names)
	holidaySet := make(map[uint16]int, len(data.Mapping))
	for _, mapping := range data.Mapping {
		index := int(mapping.Index)
		holidaySet[mapping.Day] = index
		// 名称インデックスのチェック
		if index >= namesLen {
			return nil, fmt.Errorf("invalid dataset. index overflow %d >= %d", index, namesLen)
		}
	}
	data.holidaySet = holidaySet
	return &data, nil
}

func decodeDatasetCached() func() (*StoreData, error) {
	data, err := decodeDataset()
	return func() (*StoreData, error) {
		return data, err
	}
}

// IsHoliday は指定日が祝日か返す
func IsHoliday(t time.Time) bool {
	epochDay := internal.ToEpochDay(t)
	ds, err := dataset()
	if err != nil {
		return false
	}
	set := ds.holidaySet
	_, ok := set[epochDay]
	return ok
}

// GetHolidayName は指定日の祝日名を返す
func GetHolidayName(t time.Time) *Name {
	epochDay := internal.ToEpochDay(t)
	ds, err := dataset()
	if err != nil {
		return nil
	}
	set := ds.holidaySet
	index, ok := set[epochDay]
	if !ok {
		return nil
	}
	names := ds.Names
	name := names[index]
	return &Name{Ja: name.Ja, En: name.En}
}

// Between は期間内の祝日情報を返す
func Between(start, end time.Time) []Holiday {
	epochStartDay := internal.ToEpochDay(start)
	epochEndDay := internal.ToEpochDay(end)

	ds, err := dataset()
	if err != nil {
		return nil
	}
	holidays := ds.Mapping
	names := ds.Names
	holidaysLen := len(holidays)

	// 2分探索により祝日を抽出する
	low := 0
	high := holidaysLen - 1
	startIndex := 0

	for low <= high {
		mid := (low + high) >> 1
		currentDay := holidays[mid]
		if currentDay.Day < epochStartDay {
			low = mid + 1
		} else {
			startIndex = mid
			high = mid - 1
		}
	}

	ret := make([]Holiday, 0)
	i := startIndex
	for i < holidaysLen {
		day := holidays[i]
		if day.Day > epochEndDay {
			break
		}
		name := names[day.Index]
		ret = append(ret, Holiday{
			Date: internal.FromEpochDay(day.Day),
			Name: name,
		})
		i++
	}
	return ret
}
