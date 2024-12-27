package ldate_test

import (
	"testing"
	"time"

	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/ldate"
	"github.com/stretchr/testify/assert"
)

func TestDateString(t *testing.T) {
	tt := map[string]ldate.Date{
		"2000_01_01": ldate.New(2000, 1, 1),
		"2024_12_31": ldate.New(2024, 12, 31),
		"0001_12_31": ldate.New(1, 12, 31),
		"0000_12_31": ldate.New(0, 12, 31),
		"-001_12_31": ldate.New(-1, 12, 31),
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, n, tc.String())
		})
	}
}

func TestDateValid(t *testing.T) {
	tt := []struct {
		ldate.Date
		valid bool
	}{
		{
			Date:  ldate.Date{2000, 1, 1},
			valid: true,
		},
		{
			Date:  ldate.Date{2000, 0, 1},
			valid: false,
		},
		{
			Date:  ldate.Date{2000, 1, 0},
			valid: false,
		},
		{
			Date:  ldate.Date{2000, 13, 1},
			valid: false,
		},
		{
			Date:  ldate.Date{2000, 2, 29},
			valid: true,
		},
		{
			Date:  ldate.Date{2001, 2, 29},
			valid: false,
		},
		{
			Date:  ldate.Date{2024, 4, 31},
			valid: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Date.String(), func(t *testing.T) {
			assert.Equal(t, tc.valid, tc.Date.Valid())
		})
	}
}

func TestDateNext(t *testing.T) {
	tt := []struct {
		start, next ldate.Date
	}{
		{
			start: ldate.New(2000, 1, 1),
			next:  ldate.New(2000, 1, 2),
		},
		{
			start: ldate.New(2000, 1, 31),
			next:  ldate.New(2000, 2, 1),
		},
		{
			start: ldate.New(2000, 12, 31),
			next:  ldate.New(2001, 1, 1),
		},
	}

	for _, tc := range tt {
		t.Run(tc.start.String(), func(t *testing.T) {
			assert.Equal(t, tc.next, tc.start.Next())
		})
	}
}

func TestResolve(t *testing.T) {
	tt := []struct {
		init, resolved ldate.Date
	}{
		{
			init:     ldate.New(2000, 1, 1),
			resolved: ldate.New(2000, 1, 1),
		},
		{
			init:     ldate.New(2000, 1, 0),
			resolved: ldate.New(1999, 12, 31),
		},
		{
			init:     ldate.New(2000, 1, 32),
			resolved: ldate.New(2000, 2, 1),
		},
		{
			init:     ldate.New(2000, 1, 61),
			resolved: ldate.New(2000, 3, 1),
		},
		{
			init:     ldate.New(2000, 1, 366),
			resolved: ldate.New(2000, 12, 31),
		},
		{
			init:     ldate.New(2000, 1, 367),
			resolved: ldate.New(2001, 1, 1),
		},
		{
			init:     ldate.New(2024, -10, 1),
			resolved: ldate.New(2023, 2, 1),
		},
		{
			init:     ldate.New(2024, 20, 1),
			resolved: ldate.New(2025, 8, 1),
		},
	}

	for _, tc := range tt {
		t.Run(tc.init.String(), func(t *testing.T) {
			assert.Equal(t, tc.resolved, tc.init.Resolve())
		})
	}
}

func TestSeek(t *testing.T) {
	tt := map[string]struct {
		start, expected ldate.Date
		f               filter.Filter[ldate.Date]
		max             int
		shouldFind      bool
	}{
		"03_15": {
			start: ldate.New(2024, 1, 1),
			f: func(d ldate.Date) bool {
				return d.Month == 3 && d.Day == 15
			},
			expected:   ldate.New(2024, 3, 15),
			shouldFind: true,
			max:        365,
		},
		"04_01_Wed": {
			start: ldate.New(2024, 1, 1),
			f: func(d ldate.Date) bool {
				return d.Month == 4 && d.Day == 1 && d.Weekday() == time.Wednesday
			},
			expected:   ldate.New(2026, 4, 1),
			shouldFind: true,
			max:        365 * 4,
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			got, found := tc.start.Seek(tc.f, tc.max)
			assert.Equal(t, tc.expected, got)
			assert.Equal(t, tc.shouldFind, found)
		})
	}
}

func TestWeekday(t *testing.T) {
	tt := map[string]struct {
		d        ldate.Date
		expected time.Weekday
	}{
		"2024_12_24": {
			d:        ldate.New(2024, 12, 24),
			expected: time.Tuesday,
		},
		"2024_02_03": {
			d:        ldate.New(2024, 02, 03),
			expected: time.Saturday,
		},
		"2000_07_01": {
			d:        ldate.New(2000, 7, 1),
			expected: time.Saturday,
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.d.Weekday())
		})
	}
}
