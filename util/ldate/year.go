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
