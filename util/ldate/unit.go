package ldate

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

func (u Unit) Year() Year {

	leapDays := (u / (365 * 4)) + 1

	y := (u - leapDays) / 365

	return Year(y + 2000)
}

func (u Unit) Month() Month {
	return 0
}

func (u Unit) Date() int {
	return 0
}
