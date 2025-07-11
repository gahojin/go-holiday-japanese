package internal

import (
	"time"
)

const (
	// 1日の時間 (秒)
	dayDuration = 24 * time.Hour
)

var (
	// baseTime は基準日時 (UTC)
	baseTime = time.Unix(0, 0).In(time.UTC)
)

func FromEpochDay(epochDay uint) time.Time {
	return baseTime.Add(time.Duration(epochDay) * dayDuration)
}

func ToEpochDay(date time.Time) (uint, bool) {
	// UTC時間に変換する
	year, month, day := date.Date()
	targetTime := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	if targetTime.Before(baseTime) {
		return 0, false
	}

	// 差分を計算し、日数を算出
	diff := targetTime.Sub(baseTime)
	return uint(diff / dayDuration), true
}
