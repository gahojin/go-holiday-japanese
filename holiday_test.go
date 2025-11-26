package holiday

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/gahojin/go-holiday-japanese/model"
	holidayjp "github.com/holiday-jp/holiday_jp-go"
	"github.com/stretchr/testify/assert"
)

var (
	t20250101, _ = time.Parse(time.DateOnly, "2025-01-01")
	t20250102, _ = time.Parse(time.DateOnly, "2025-01-02")
	t20250111, _ = time.Parse(time.DateOnly, "2025-01-11")
	t20250112, _ = time.Parse(time.DateOnly, "2025-01-12")
	t20250113, _ = time.Parse(time.DateOnly, "2025-01-13")
	t20250131, _ = time.Parse(time.DateOnly, "2025-01-31")
	t20240230, _ = time.Parse(time.DateOnly, "2024-02-30")
	t20240301, _ = time.Parse(time.DateOnly, "2024-03-01")
	t20241014, _ = time.Parse(time.DateOnly, "2024-10-14")
	t20241015, _ = time.Parse(time.DateOnly, "2024-10-15")
	t20200101, _ = time.Parse(time.DateOnly, "2020-01-01")
	t20201231, _ = time.Parse(time.DateOnly, "2020-12-31")
	t19700101, _ = time.Parse(time.DateOnly, "1970-01-01")
	t19700115, _ = time.Parse(time.DateOnly, "1970-01-15")

	tokyoLoc, _      = time.LoadLocation("Asia/Tokyo")
	newYorkLoc, _    = time.LoadLocation("America/New_York")
	losAngelesLoc, _ = time.LoadLocation("America/Los_Angeles")
	kolkataLoc, _    = time.LoadLocation("Asia/Kolkata")
)

func BenchmarkIsHoliday(b *testing.B) {
	IsHoliday(t20241014)
	IsHoliday(t20241015)
}

func BenchmarkIsHolidayHolidayJp(b *testing.B) {
	holidayjp.IsHoliday(t20241014)
	holidayjp.IsHoliday(t20241015)
}

func BenchmarkBetween(b *testing.B) {
	Between(t20200101, t20201231)
}

func BenchmarkBetweenHolidayJp(b *testing.B) {
	holidayjp.Between(t20200101, t20201231)
}

func TestIsHoliday(t *testing.T) {
	tests := []struct {
		name string
		tz   *time.Location
	}{
		{
			name: "JST (UTC+9)",
			tz:   tokyoLoc,
		},
		{
			name: "EST (UTC-5)",
			tz:   newYorkLoc, // 東部標準時 (UTC-5)
		},
		{
			name: "PST (UTC-8)",
			tz:   losAngelesLoc, // 太平洋標準時 (UTC-8)
		},
		{
			name: "IST (UTC+05:30)",
			tz:   kolkataLoc, // UTC+05:30
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s: 特定の日が祝日か", tt.name), func(t *testing.T) {
			// スポーツの日
			assert.True(t, IsHoliday(time.Date(2024, 10, 14, 0, 0, 0, 0, tt.tz)))
			assert.False(t, IsHoliday(time.Date(2024, 10, 13, 0, 0, 0, 0, tt.tz)))
			assert.False(t, IsHoliday(time.Date(2024, 10, 15, 0, 0, 0, 0, tt.tz)))

			// 山の日
			assert.False(t, IsHoliday(time.Date(2015, 8, 11, 0, 0, 0, 0, tt.tz)))
			for year := 2016; year <= 2050; year++ {
				switch year {
				case 2020:
					assert.True(t, IsHoliday(time.Date(year, 8, 10, 0, 0, 0, 0, tt.tz)))
					assert.False(t, IsHoliday(time.Date(year, 8, 11, 0, 0, 0, 0, tt.tz)))
				case 2021:
					assert.True(t, IsHoliday(time.Date(year, 8, 8, 0, 0, 0, 0, tt.tz)))
					assert.False(t, IsHoliday(time.Date(year, 8, 11, 0, 0, 0, 0, tt.tz)))
				default:
					assert.True(t, IsHoliday(time.Date(year, 8, 11, 0, 0, 0, 0, tt.tz)))
				}
			}
			// 山の日 (振替休日)
			assert.True(t, IsHoliday(time.Date(2021, 8, 9, 0, 0, 0, 0, tt.tz)))
		})
	}
}

func TestGetHolidayName(t *testing.T) {
	tests := []struct {
		name      string
		year      int
		month     int
		day       int
		holidayJa string
	}{
		{
			year:      2025,
			month:     1,
			day:       1,
			holidayJa: "元日",
		},
		{
			year:      2025,
			month:     1,
			day:       13,
			holidayJa: "成人の日",
		},
		{
			year:      2024,
			month:     10,
			day:       14,
			holidayJa: "スポーツの日",
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s: 祝日名があっていること", tt.name), func(t *testing.T) {
			assert.Equal(t, GetHolidayName(time.Date(tt.year, time.Month(tt.month), tt.day, 0, 0, 0, 0, time.UTC)).Ja, tt.holidayJa)
		})
	}
}

func TestBetween(t *testing.T) {
	// 境界チェック
	assert.Empty(t, Between(t20250102, t20250112))
	assert.Equal(t, []model.Holiday{
		{Date: t20250101, Name: model.Name{Ja: "元日", En: "New Year's Day"}},
	}, Between(t20250101, t20250112))
	assert.Equal(t, []model.Holiday{
		{Date: t20250113, Name: model.Name{Ja: "成人の日", En: "Coming of Age Day"}},
	}, Between(t20250113, t20250131))
	assert.Equal(t, []model.Holiday{
		{Date: t20250113, Name: model.Name{Ja: "成人の日", En: "Coming of Age Day"}},
	}, Between(t20250111, t20250131))
	assert.Equal(t, []model.Holiday{
		{Date: t20250101, Name: model.Name{Ja: "元日", En: "New Year's Day"}},
	}, Between(t20250101, t20250101))
	assert.Equal(t, []model.Holiday{
		{Date: t19700101, Name: model.Name{Ja: "元日", En: "New Year's Day"}},
		{Date: t19700115, Name: model.Name{Ja: "成人の日", En: "Coming of Age Day"}},
	}, Between(t19700101, t19700115))
	assert.Empty(t, Between(t20240230, t20240301))
}

func TestCompatibility(t *testing.T) {
	year := 1970
	for year < 2025 {
		s := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
		e := time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC)
		thisLibrary := Between(s, e)
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
