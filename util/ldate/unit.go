package ldate

// == projects.Code.luce.ldate ==
// [ ] rigorous years
//	Unit.Year, Unit.Leapdays and Year.Leapdays do not match the calendar
//	(TestUnitYear, TestLeapDays and TestYearLeapDays fail), and Unit.Month and
//	Unit.Date are not written yet. The skipped tests wait on this item.

// Unit is a number of days, counted from 2000_01_01 as Unit 0. It is unfinished,
// see the project marker above.
type Unit int64

const (
	march28   = 89
	oneYear   = 365
	fourYears = oneYear * 4
)

// Leapdays between u and and Unit(0). Note that if u is negative, the number
// of leap days will be negative.
func (u Unit) Leapdays() int64 {
	// x is relative to 1996_03_29
	x := int64(u) - march28 + fourYears + 1
	ld := int64(x) / fourYears
	x += ld
	for {
		newLD := int64(x) / fourYears
		if ld == newLD {
			break
		}
		ld = newLD
	}
	return int64(x) / fourYears
}

// Year returns the Year of the Unit. It does not match the calendar for every
// Unit yet.
func (u Unit) Year() Year {

	leapDays := (u / (365 * 4)) + 1

	y := (u - leapDays) / 365

	return Year(y + 2000)
}

// Month is not written yet and returns 0.
func (u Unit) Month() Month {
	return 0
}

// Date is not written yet and returns 0.
func (u Unit) Date() int {
	return 0
}
