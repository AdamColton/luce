package ldate_test

import (
	"sort"
	"testing"
	"time"

	"github.com/adamcolton/luce/util/ldate"
	"github.com/stretchr/testify/assert"
)

func TestLess(t *testing.T) {
	dates := []ldate.Date{
		ldate.New(2024, 6, 15),
		ldate.New(2023, 12, 31),
		ldate.New(2024, 6, 14),
		ldate.New(2024, 1, 1),
	}
	sort.Slice(dates, func(i, j int) bool { return ldate.Less(dates[i], dates[j]) })
	assert.Equal(t, []ldate.Date{
		ldate.New(2023, 12, 31),
		ldate.New(2024, 1, 1),
		ldate.New(2024, 6, 14),
		ldate.New(2024, 6, 15),
	}, dates)
	assert.False(t, ldate.Less(dates[0], dates[0]))
	assert.True(t, ldate.Less(dates[0], dates[1]) == dates[0].Before(dates[1]))
}

func TestFromTime(t *testing.T) {
	afternoon := time.Date(2024, time.February, 29, 15, 30, 0, 0, time.UTC)
	assert.Equal(t, ldate.New(2024, 2, 29), ldate.FromTime(afternoon))
}
