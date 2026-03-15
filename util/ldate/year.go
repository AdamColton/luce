package ldate

import (
	"fmt"
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

// Leapdays is unfinished: it counts every fourth year since 1997 and ignores the
// century rules. See the project marker in unit.go.
func (y Year) Leapdays() int64 {
	y64 := int64(y)
	return ((y64 - 1997) / 4) //- ((y64 - 2001) / 100) + ((y64 - 2001) / 400)
}
