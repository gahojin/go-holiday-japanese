package internal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFromEpochDay(t *testing.T) {
	tests := []struct {
		epochDay uint16
		expected time.Time
	}{
		{0, time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)},
		{1, time.Date(1970, 1, 2, 0, 0, 0, 0, time.UTC)},
		{20089, time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.expected, FromEpochDay(tt.epochDay))
	}
}

func TestToEpochDay(t *testing.T) {
	tests := []struct {
		date     time.Time
		expected uint16
		ok       bool
	}{
		{time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), 0, true},
		{time.Date(1970, 1, 1, 23, 59, 59, 0, time.UTC), 0, true},
		{time.Date(1970, 1, 2, 0, 0, 0, 0, time.UTC), 1, true},
		{time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC), 20089, true},
		{time.Date(1969, 12, 31, 23, 59, 59, 0, time.UTC), 0, false},
		{time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC), 0, false},
	}
	for _, tt := range tests {
		day, ok := ToEpochDay(tt.date)
		assert.Equal(t, tt.ok, ok)
		if tt.ok {
			assert.Equal(t, tt.expected, day)
		}
	}
}

func TestConvertDataset(t *testing.T) {
	// \x01\x00 -> day += 1, index = 0
	// \x02\x02 -> day += 2, index = 2 (day = 1+2=3)
	mapping := "\x01\x00\x02\x02"
	bitset, results := ConvertDataset(mapping)

	assert.True(t, bitset.Has(1))
	assert.True(t, bitset.Has(3))

	assert.Len(t, results, 2)
	assert.Equal(t, uint16(1), results[0].Day)
	assert.Equal(t, uint8(0), results[0].Index)
	assert.Equal(t, uint16(3), results[1].Day)
	assert.Equal(t, uint8(2), results[1].Index)
}
