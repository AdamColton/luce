package ldate

import (
	"fmt"

	"github.com/adamcolton/luce/math/ints"
)

// Year is a year of the Gregorian calendar.
type Year int

// IsLeapYear uses the Gregorian rule: every fourth year, except the centuries
// that are not divisible by 400.
func (y Year) IsLeapYear() bool {
	return (y%400 == 0) || (y%4 == 0 && y%100 != 0)
}

// String formats the year with at least four characters: 0007, 2024, -001.
func (y Year) String() string {
	return fmt.Sprintf("%04d", y)
}

// Leapdays returns the number of leap years from 2000 up to, but not including,
// y. For a year before 2000 it is negative: minus the number from y up to 2000.
// It matches Ordinal.Leapdays for January 1st of y.
func (y Year) Leapdays() int64 {
	return leapYearsBefore(int64(y)) - leapYearsBefore(2000)
}

// leapYearsBefore counts the leap years in [0, y), negative when y < 0.
// Year 0 is a leap year, which is the +1.
func leapYearsBefore(y int64) int64 {
	y--
	return ints.DivDown(y, 4) - ints.DivDown(y, 100) + ints.DivDown(y, 400) + 1
}
