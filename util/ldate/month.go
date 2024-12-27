package ldate

import (
	"golang.org/x/exp/constraints"
)

// Month is a month of the year, from 1 (January) to 12 (December). Other values
// are not valid: they have no name and no days.
type Month byte

// The months of the year.
const (
	January Month = 1 + iota
	February
	March
	April
	May
	June
	July
	August
	September
	October
	November
	December
)

// Short names for the months. May has none because it is short already.
const (
	Jan = January
	Feb = February
	Mar = March
	Apr = April
	Jun = June
	Jul = July
	Aug = August
	Sep = September
	Oct = October
	Nov = November
	Dec = December
)

var (
	names = [...]string{
		"none",
		"January",
		"February",
		"March",
		"April",
		"May",
		"June",
		"July",
		"August",
		"September",
		"October",
		"November",
		"December",
	}

	shortNames = [...]string{
		"none",
		"Jan",
		"Feb",
		"Mar",
		"Apr",
		"May",
		"Jun",
		"Jul",
		"Aug",
		"Sep",
		"Oct",
		"Nov",
		"Dec",
	}

	intTab = [...]string{
		"00",
		"01",
		"02",
		"03",
		"04",
		"05",
		"06",
		"07",
		"08",
		"09",
		"10",
		"11",
		"12",
		"13",
		"14",
		"15",
		"16",
		"17",
		"18",
		"19",
		"20",
		"21",
		"22",
		"23",
		"24",
		"25",
		"26",
		"27",
		"28",
		"29",
		"30",
		"31",
	}

	daysInMonthTab = [...]int{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
)

// Name returns the name of the month, or its three letter abbreviation if short
// is true. It returns "none" for a Month that is out of range.
func (m Month) Name(short bool) string {
	if m > 12 {
		return names[0]
	}
	if short {
		return shortNames[m]
	}
	return names[m]
}

// String returns the month as two digits, 01 to 12, and 00 for a Month that is
// out of range.
func (m Month) String() string {
	if m > 12 {
		return intTab[0]
	}
	return intTab[m]
}

// Valid returns true for 1 to 12.
func (m Month) Valid() bool {
	return m > 0 && m <= 12
}

// Days returns the number of days in the month in year y, and 0 for a Month
// that is out of range.
func (m Month) Days(y Year) int {
	if m == Feb && y.IsLeapYear() {
		return 29
	}
	if m > 12 {
		return 0
	}
	return daysInMonthTab[m]
}

// MonthYear returns the Month and Year for month m of year y, carrying a month
// out of range into the year: MonthYear(14, 2024) is February 2025 and
// MonthYear(0, 2024) is December 2023.
func MonthYear[M, Y constraints.Integer](m M, y Y) (Month, Year) {
	yy := Year(y)
	if m > 0 && m < 13 {
		return Month(m), yy
	}
	for m < 1 {
		yy--
		m += 12
	}
	for m > 12 {
		yy++
		m -= 12
	}
	return Month(m), yy
}
