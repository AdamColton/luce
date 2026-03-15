package ldate

import (
	"fmt"
)

type Year int

func (y Year) IsLeapYear() bool {
	//return y%4 == 0 && (y%100 != 0 || y%500 == 0)

	return (y%400 == 0) || (y%4 == 0 && y%100 != 0)
}

func (y Year) String() string {
	return fmt.Sprintf("%04d", y)
}

func (y Year) Leapdays() int64 {
	y64 := int64(y)
	return ((y64 - 1997) / 4) //- ((y64 - 2001) / 100) + ((y64 - 2001) / 400)
}
