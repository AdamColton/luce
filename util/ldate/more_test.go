package ldate_test

import (
	"testing"
	"time"

	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/ldate"
	"github.com/stretchr/testify/assert"
)

func TestNow(t *testing.T) {
	today := ldate.TimeToDate(time.Now())
	justBefore := ldate.TimeToDate(time.Now().Add(-time.Minute))
	assert.Contains(t, []ldate.Date{today, justBefore}, ldate.Now())
}

func TestBefore(t *testing.T) {
	d := ldate.New(2024, 6, 15)
	tt := map[string]struct {
		other  ldate.Date
		before bool
	}{
		"later year":    {ldate.New(2025, 1, 1), true},
		"earlier year":  {ldate.New(2023, 12, 31), false},
		"later month":   {ldate.New(2024, 7, 1), true},
		"earlier month": {ldate.New(2024, 5, 31), false},
		"later day":     {ldate.New(2024, 6, 16), true},
		"earlier day":   {ldate.New(2024, 6, 14), false},
		"same":          {d, false},
	}
	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.before, d.Before(tc.other))
		})
	}
}

func TestSeekLimit(t *testing.T) {
	start := ldate.New(2024, 2, 27)
	firstOfMonth := filter.New(func(d ldate.Date) bool { return d.Day == 1 })

	d, ok := start.Seek(firstOfMonth, 3)
	assert.True(t, ok, "the last day it looks at can be the one")
	assert.Equal(t, ldate.New(2024, 3, 1), d)

	d, ok = start.Seek(firstOfMonth, 2)
	assert.False(t, ok)
	assert.Equal(t, ldate.New(2024, 2, 29), d, "the last day it looked at")
}

func TestMonthNames(t *testing.T) {
	assert.Equal(t, "March", ldate.March.Name(false))
	assert.Equal(t, "Mar", ldate.March.Name(true))
	assert.Equal(t, "none", ldate.Month(13).Name(false))
	assert.Equal(t, "none", ldate.Month(0).Name(true))
	assert.Equal(t, "03", ldate.March.String())
	assert.Equal(t, "00", ldate.Month(13).String())
}

func TestMonthDays(t *testing.T) {
	assert.Equal(t, 31, ldate.January.Days(2023))
	assert.Equal(t, 28, ldate.February.Days(2023))
	assert.Equal(t, 29, ldate.February.Days(2024))
	assert.Equal(t, 0, ldate.Month(13).Days(2024), "no days in a month that does not exist")
}

func TestMonthYear(t *testing.T) {
	m, y := ldate.MonthYear(14, 2024)
	assert.Equal(t, ldate.February, m)
	assert.Equal(t, ldate.Year(2025), y)

	m, y = ldate.MonthYear(0, 2024)
	assert.Equal(t, ldate.December, m)
	assert.Equal(t, ldate.Year(2023), y)
}

func TestYearString(t *testing.T) {
	assert.Equal(t, "0007", ldate.Year(7).String())
	assert.Equal(t, "2024", ldate.Year(2024).String())
	assert.Equal(t, "-001", ldate.Year(-1).String())
}
