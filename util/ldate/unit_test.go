package ldate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnitYear(t *testing.T) {
	cur := time.Date(2000, time.January, 1, 12, 0, 0, 0, time.UTC)
	for u := range Unit(8000) {
		if !assert.Equal(t, Year(cur.Year()), u.Year(), "%v %d", cur, u) {
			break
		}

		cur = cur.Add(24 * time.Hour)
	}
}

func TestLeapDays(t *testing.T) {
	cur := time.Date(2000, time.January, 1, 12, 0, 0, 0, time.UTC)
	var leapdays int64
	for u := range Unit(8000) {
		if !assert.Equal(t, leapdays, u.Leapdays(), "%v %d", cur, u) {
			break
		}

		cur = cur.Add(24 * time.Hour)
		if cur.Month() == time.March && cur.Day() == 29 {
			leapdays++
		}
	}
}

func TestYearLeapDays(t *testing.T) {
	var ld int64
	for y := Year(2000); y < 4000; y++ {
		if !assert.Equal(t, ld, y.Leapdays(), y) {
			break
		}
		if y.IsLeapYear() {
			ld++
		}
	}

	ld = 0
	for y := Year(1999); y > 0; y-- {
		if y.IsLeapYear() {
			ld--
		}
		if !assert.Equal(t, ld, y.Leapdays(), y) {
			break
		}
	}
}
