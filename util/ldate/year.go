package ldate

import (
	"fmt"
)

type Year int

func (y Year) IsLeapYear() bool {
	return y%4 == 0 && (y%100 != 0 || y%500 == 0)
}

func (y Year) String() string {
	return fmt.Sprintf("%04d", y)
}
