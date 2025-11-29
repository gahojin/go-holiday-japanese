package compatible

import (
	"slices"
	"testing"
	"time"

	"github.com/gahojin/go-holiday-japanese"
	"github.com/gahojin/go-holiday-japanese/model"
	holidayjp "github.com/holiday-jp/holiday_jp-go"
	"github.com/stretchr/testify/assert"
)

var (
	t20241014, _ = time.Parse(time.DateOnly, "2024-10-14")
	t20241015, _ = time.Parse(time.DateOnly, "2024-10-15")
	t20200101, _ = time.Parse(time.DateOnly, "2020-01-01")
	t20201231, _ = time.Parse(time.DateOnly, "2020-12-31")
)

func BenchmarkIsHoliday(b *testing.B) {
	for i := 0; i < b.N; i++ {
		holiday.IsHoliday(t20241014)
		holiday.IsHoliday(t20241015)
	}
}

func BenchmarkIsHolidayHolidayJp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		holidayjp.IsHoliday(t20241014)
		holidayjp.IsHoliday(t20241015)
	}
}

func BenchmarkBetween(b *testing.B) {
	for i := 0; i < b.N; i++ {
		holiday.Between(t20200101, t20201231)
	}
}

func BenchmarkBetweenHolidayJp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		holidayjp.Between(t20200101, t20201231)
	}
}

func TestCompatibility(t *testing.T) {
	year := 1970
	for year < 2025 {
		s := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
		e := time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC)
		thisLibrary := holiday.Between(s, e)
		holidayJp := holidayjp.Between(s, e)

		// 比較
		holidayJpList := make([]model.Holiday, len(thisLibrary))
		index := 0
		for key, value := range holidayJp {
			// holiday-jpのHoliday.Date関数は動作しないため、keyから日付パース
			date, err := time.Parse(time.DateOnly, key)
			assert.NoError(t, err)
			holidayJpList[index] = model.Holiday{
				Date: date,
				Name: model.Name{
					Ja: value.Name(),
					En: value.NameEn(),
				},
			}
			index++
		}
		// mapは日付順ではないため、日付順にソート
		slices.SortFunc(holidayJpList, func(a, b model.Holiday) int {
			return a.Date.Compare(b.Date)
		})

		assert.Equal(t, holidayJpList, thisLibrary)

		year++
	}
}
