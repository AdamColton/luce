package ldate_test

import (
	"strconv"
	"testing"

	"github.com/adamcolton/luce/util/ldate"
	"github.com/stretchr/testify/assert"
)

func TestIsLeapYear(t *testing.T) {
	tt := map[ldate.Year]bool{
		2000: true,
		2001: false,
		2002: false,
		2003: false,
		2004: true,
		2100: false,
	}

	for y, expected := range tt {
		t.Run(strconv.Itoa(int(y)), func(t *testing.T) {
			assert.Equal(t, expected, y.IsLeapYear())
		})
	}
}
