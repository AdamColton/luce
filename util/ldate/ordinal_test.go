package ldate_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/adamcolton/luce/util/ldate"
	"github.com/stretchr/testify/assert"
)

func ExampleOrdinal() {
	moonLanding := ldate.New(1969, 7, 20).Ordinal()
	millennium := ldate.New(2000, 1, 1).Ordinal()
	fmt.Println(moonLanding, millennium)
	fmt.Println(millennium-moonLanding, "days apart")
	fmt.Println((moonLanding + 10000).Date())
	// Output:
	// -11122 0
	// 11122 days apart
	// 1996_12_05
}

// TestOrdinalDate walks day by day across more than 800 years, through the
// century years 1700, 1800, 1900 and 2100 that are not leap years and 1600,
// 2000 and 2400 that are, checking against time.Time.
func TestOrdinalDate(t *testing.T) {
	const span = 150_000
	start := time.Date(2000, time.January, 1-span, 12, 0, 0, 0, time.UTC)
	for o := ldate.Ordinal(-span); o <= span; o++ {
		cur := start.AddDate(0, 0, int(o+span))
		expected := ldate.FromTime(cur)
		if !assert.Equal(t, expected, o.Date(), "%v %d", cur, o) ||
			!assert.Equal(t, o, expected.Ordinal(), "%v %d", cur, o) {
			break
		}
	}
}

// TestOrdinalDateFar samples Ordinals far from 2000 in both directions, through year
// 0 and into negative years.
func TestOrdinalDateFar(t *testing.T) {
	const span, step = 100_000_000, 9_973
	millennium := time.Date(2000, time.January, 1, 12, 0, 0, 0, time.UTC)
	for o := ldate.Ordinal(-span); o <= span; o += step {
		cur := millennium.AddDate(0, 0, int(o))
		expected := ldate.FromTime(cur)
		if !assert.Equal(t, expected, o.Date(), "%v %d", cur, o) ||
			!assert.Equal(t, o, expected.Ordinal(), "%v %d", cur, o) {
			break
		}
	}

	yearZero := ldate.New(0, 1, 1).Ordinal()
	assert.Equal(t, ldate.Ordinal(-730485), yearZero)
	assert.Equal(t, ldate.New(-1, 12, 31), (yearZero - 1).Date())
	assert.Equal(t, ldate.New(0, 2, 29), (yearZero + 59).Date(), "year 0 is a leap year")
}

func TestOrdinalFields(t *testing.T) {
	leapDay := ldate.New(2024, 2, 29).Ordinal()
	assert.Equal(t, ldate.Year(2024), leapDay.Year())
	assert.Equal(t, ldate.February, leapDay.Month())
	assert.Equal(t, 29, leapDay.Day())
	assert.Equal(t, ldate.New(2024, 3, 1), (leapDay + 1).Date())
}

func TestDateOrdinalCarries(t *testing.T) {
	assert.Equal(t, ldate.New(2024, 3, 1).Ordinal(), ldate.Date{Year: 2024, Month: 2, Day: 30}.Ordinal())
	assert.Equal(t, ldate.New(2023, 12, 31).Ordinal(), ldate.Date{Year: 2024, Month: 1, Day: 0}.Ordinal())
	assert.Equal(t, ldate.New(2025, 2, 1).Ordinal(), ldate.Date{Year: 2024, Month: 14, Day: 1}.Ordinal())
}
