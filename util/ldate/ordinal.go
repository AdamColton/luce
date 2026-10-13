package ldate

import "github.com/adamcolton/luce/math/ints"

// Ordinal is a number of days, counted from 2000_01_01 as Ordinal 0. Days
// before that are negative. An Ordinal converts to and from a Date in constant
// time, and differences between Ordinals are differences in days.
type Ordinal int64

// The conversions follow Howard Hinnant's days_from_civil and civil_from_days
// (http://howardhinnant.github.io/date_algorithms.html). The Gregorian calendar
// repeats every 400 years (an era), and counting each era from March 1st puts
// the leap day at the end of the year, so the month and day can be computed
// without knowing whether the year is a leap year.
const (
	daysPerEra = 146097 // 400*365 + 97 leap days
	// march1st2000 is the Ordinal of 2000_03_01, where the eras start.
	march1st2000 = 31 + 29
)

// Date returns the calendar date of the Ordinal.
func (o Ordinal) Date() Date {
	z := int64(o) - march1st2000
	// Rounding down puts days before 2000_03_01 in the eras before it.
	era := ints.DivDown(z, daysPerEra)
	doe := ints.Mod(z, daysPerEra) // day of era, [0, 146096]
	// Remove the leap days added every 4 years (1460 days), the ones missing every
	// 100 years (36524 days) and the one added back every 400 years (146096 days,
	// only the last day of the era) to get the year of the era, [0, 399].
	yoe := (doe - doe/1460 + doe/36524 - doe/146096) / 365
	doy := doe - (365*yoe + yoe/4 - yoe/100) // day of year from March 1st, [0, 365]
	// Months from March: March to July and August to December each run
	// 31,30,31,30,31 days, 153 days in all, so a month is 153/5 days on average.
	mp := (5*doy + 2) / 153 // [0, 11]
	d := doy - (153*mp+2)/5 + 1
	m := mp + 3
	y := 2000 + era*400 + yoe
	if mp >= 10 {
		// January and February belong to the next calendar year.
		m -= 12
		y++
	}
	return Date{
		Year:  Year(y),
		Month: Month(m),
		Day:   int(d),
	}
}

// Year returns the year of the Ordinal.
func (o Ordinal) Year() Year {
	return o.Date().Year
}

// Month returns the month of the Ordinal.
func (o Ordinal) Month() Month {
	return o.Date().Month
}

// Day returns the day of the month of the Ordinal.
func (o Ordinal) Day() int {
	return o.Date().Day
}

// Ordinal returns the number of days from 2000_01_01 to d. A Day out of range
// carries into the neighboring months, as in New.
func (d Date) Ordinal() Ordinal {
	m, y := MonthYear(d.Month, d.Year)
	// Count years from March 1st, see Ordinal.Date.
	if m <= 2 {
		y--
	}
	ys := int64(y) - 2000
	era := ints.DivDown(ys, 400)
	yoe := ints.Mod(ys, 400)
	mp := (int64(m) + 9) % 12
	doy := (153*mp+2)/5 + int64(d.Day) - 1
	doe := yoe*365 + yoe/4 - yoe/100 + doy
	return Ordinal(era*daysPerEra + doe + march1st2000)
}

// Leapdays returns the number of February 29ths from Ordinal(0) up to, but not
// including, o. For a negative o it is negative: minus the number from o up to
// Ordinal(0).
func (o Ordinal) Leapdays() int64 {
	d := o.Date()
	ld := d.Year.Leapdays()
	if d.Month > February && d.Year.IsLeapYear() {
		ld++
	}
	return ld
}
