package ldate

import (
	"strconv"
	"time"

	"github.com/adamcolton/luce/util/filter"
)

// Date is a day in the calendar. A Date is not normalized: New and Resolve turn
// a Month or Day that is out of range into the date it means.
type Date struct {
	Year  Year
	Month Month
	Day   int
}

// New returns the Date for day d of month m of year y. A month or day out of
// range carries into the neighboring months and years, so New(2024, 13, 1) is
// January 1st 2025 and New(2024, 3, 0) is the last day of February.
func New(y, m, d int) Date {
	mm, yy := MonthYear(m, y)
	for {
		if d < 1 {
			mm, yy = MonthYear(mm-1, yy)
			d += mm.Days(yy)
			continue
		}
		md := mm.Days(yy)
		if d > md {
			mm, yy = MonthYear(mm+1, yy)
			d -= md
			continue
		}
		break
	}

	return Date{
		Year:  yy,
		Month: mm,
		Day:   d,
	}
}

// Valid returns true if the Month is 1 to 12 and the Day exists in it.
func (d Date) Valid() bool {
	return d.Month.Valid() && d.Day > 0 && d.Day <= d.Month.Days(d.Year)
}

// String formats the Date as YYYY_MM_DD. The year is padded to four characters
// (a negative year keeps its sign: -001) and a Month or Day that is out of
// range is written as 00.
func (d Date) String() string {
	out := make([]byte, 0, 12)
	out = strconv.AppendInt(out, int64(d.Year), 10)
	ln := len(out)
	if pad := 4 - ln; pad > 0 {
		out = out[:4]
		copy(out[4-ln:], out)
		for i := 0; i < pad; i++ {
			out[i] = '0'
		}
		if d.Year < 0 {
			signIdx := 4 - ln
			out[0], out[signIdx] = out[signIdx], out[0]
		}
	}

	out = append(out, '_')
	out = append(out, d.Month.String()...)
	out = append(out, '_')
	dIdx := d.Day
	if d.Day < 1 || d.Day > 31 {
		dIdx = 0
	}
	out = append(out, intTab[dIdx]...)
	return string(out)
}

// Next returns the day after d. d must be valid.
func (d Date) Next() Date {
	d.Day++
	if d.Day > d.Month.Days(d.Year) {
		d.Day = 1
		d.Month, d.Year = MonthYear(d.Month+1, d.Year)
	}
	return d
}

// Resolve returns the normalized Date for d, see New.
func (d Date) Resolve() Date {
	return New(int(d.Year), int(d.Month), d.Day)
}

// Seek looks for the first Date, starting with d and going forward, for which
// f is true. It looks at d and at most max days after it. If it finds one it
// returns it and true, otherwise the last Date it looked at and false.
func (d Date) Seek(f filter.Filter[Date], max int) (Date, bool) {
	for ; max > 0; max-- {
		if f(d) {
			return d, true
		}
		d = d.Next()
	}
	return d, f(d)
}

var monthKey = []int{0, 3, 2, 5, 0, 3, 5, 1, 4, 6, 2, 4}

// Weekday returns the day of the week. The Month must be 1 to 12, otherwise
// Weekday panics.
func (d Date) Weekday() time.Weekday {
	y := int(d.Year)
	m := int(d.Month)
	dd := int(d.Day)
	if m < 3 {
		y--
	}
	return time.Weekday((y + y/4 - y/100 + y/400 + monthKey[m-1] + dd) % 7)
}

// Before returns true if d is before d2.
func (d Date) Before(d2 Date) bool {
	if d.Year < d2.Year {
		return true
	}
	if d.Year > d2.Year {
		return false
	}
	if d.Month < d2.Month {
		return true
	}
	if d.Month > d2.Month {
		return false
	}
	if d.Day < d2.Day {
		return true
	}
	if d.Day > d2.Day {
		return false
	}
	return false
}

// TimeToDate returns the Date of t, in t's location.
func TimeToDate(t time.Time) Date {
	return Date{
		Year:  Year(t.Year()),
		Month: Month(t.Month()),
		Day:   t.Day(),
	}
}
