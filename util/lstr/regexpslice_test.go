package lstr_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/util/lstr"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

func TestNewNamedRegexpSlice(t *testing.T) {
	number := `^\s*(\d*\.?\d+)\s*$`
	word := `^\s*([a-zA-Z]+)\s*$`
	unit := `^\s*(\d*\.?\d+) +([a-zA-Z]+)\s*$`
	rm := lstr.NewRegexpSlice(
		number,
		word,
		unit,
	)
	str := `

		3.1415
		apple
		not a match
		5 cups

		aardvark
		
	`
	strs := lstr.NewLine.Strings(str)
	expected := slice.New([][]string{
		{number, "3.1415"},
		{word, "apple"},
		{unit, "5", "cups"},
		{word, "aardvark"},
	}).Iter()

	err := timeout.After(100000000, func() {
		for !strs.Done() {
			id, m := strs.RegMap(rm)
			exp, _ := expected.Cur()
			expected.Next()
			assert.Equal(t, exp[0], id)
			assert.Equal(t, exp[1:], m[1:])
		}
	})
	assert.NoError(t, err)
	assert.True(t, expected.Done())

	id, m := strs.RegMap(rm)
	assert.Equal(t, "", id)
	assert.Nil(t, m)
}
