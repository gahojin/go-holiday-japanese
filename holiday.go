package holiday

import (
	_ "embed"
	"sort"
	"time"

	"github.com/gahojin/go-holiday-japanese/internal"
	"github.com/gahojin/go-holiday-japanese/model"
)

var holidayBitset, mapping = internal.ConvertDataset(internal.HolidayMapping)

// IsHoliday は指定日が祝日か返す
func IsHoliday(t time.Time) bool {
	epochDay, ok := internal.ToEpochDay(t)
	if !ok {
		return false
	}
	return holidayBitset.Has(epochDay)
}

// GetHolidayName は指定日の祝日名を返す
func GetHolidayName(t time.Time) *model.Name {
	epochDay, ok := internal.ToEpochDay(t)
	if !ok {
		return nil
	}

	mappingLen := len(mapping)
	idx := sort.Search(mappingLen, func(i int) bool {
		return mapping[i].Day >= epochDay
	})
	if idx >= mappingLen {
		return nil
	}
	data := mapping[idx]
	if data.Day != epochDay {
		return nil
	}
	index := data.Index
	names := internal.HolidayNames
	return &model.Name{Ja: internal.HolidayNames[index], En: names[index+1]}
}

// Between は期間内の祝日情報を返す
func Between(start, end time.Time) []model.Holiday {
	epochStartDay, ok := internal.ToEpochDay(start)
	if !ok {
		return nil
	}
	epochEndDay, ok := internal.ToEpochDay(end)
	if !ok {
		return nil
	}

	names := internal.HolidayNames
	mappingLen := len(mapping)

	// 2分探索により祝日を抽出する
	startIndex := sort.Search(mappingLen, func(i int) bool {
		return mapping[i].Day >= epochStartDay
	})
	endIndex := sort.Search(mappingLen, func(i int) bool {
		return mapping[i].Day > epochEndDay
	})

	// あらかじめ確保するサイズを算出し，メモリ最適化する
	count := endIndex - startIndex
	if count <= 0 {
		return []model.Holiday{}
	}
	ret := make([]model.Holiday, count)

	j := 0
	for i := startIndex; i < endIndex; i++ {
		day := mapping[i]
		ret[j] = model.Holiday{
			Date: internal.FromEpochDay(day.Day),
			Name: model.Name{
				Ja: names[day.Index],
				En: names[day.Index+1],
			},
		}
		j++
	}
	return ret
}
