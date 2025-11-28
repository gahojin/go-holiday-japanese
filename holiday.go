package holiday

import (
	_ "embed"
	"time"

	"github.com/gahojin/go-holiday-japanese/internal"
	"github.com/gahojin/go-holiday-japanese/model"
)

var holidays, mapping = internal.ConvertDataset(holidayMapping)

// IsHoliday は指定日が祝日か返す
func IsHoliday(t time.Time) bool {
	epochDay, ok := internal.ToEpochDay(t)
	if !ok {
		return false
	}
	_, ok = holidays[epochDay]
	return ok
}

// GetHolidayName は指定日の祝日名を返す
func GetHolidayName(t time.Time) *model.Name {
	epochDay, ok := internal.ToEpochDay(t)
	if !ok {
		return nil
	}
	index, ok := holidays[epochDay]
	if !ok {
		return nil
	}
	names := holidayNames
	return &model.Name{Ja: names[index], En: names[index+1]}
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

	names := holidayNames
	mappingLen := len(mapping)

	// 2分探索により祝日を抽出する
	low := 0
	high := mappingLen - 1
	startIndex := high + 1

	for low <= high {
		mid := (low + high) >> 1
		currentDay := mapping[mid]
		if currentDay.Day < epochStartDay {
			low = mid + 1
		} else {
			startIndex = mid
			high = mid - 1
		}
	}

	ret := make([]model.Holiday, 0)
	i := startIndex
	for i < mappingLen {
		day := mapping[i]
		if day.Day > epochEndDay {
			break
		}
		ret = append(ret, model.Holiday{
			Date: internal.FromEpochDay(day.Day),
			Name: model.Name{
				Ja: names[day.Index],
				En: names[day.Index+1],
			},
		})
		i++
	}
	return ret
}
