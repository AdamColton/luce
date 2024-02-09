package lstr_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/util/liter"
	"github.com/adamcolton/luce/util/lstr"
	"github.com/stretchr/testify/assert"
)

func TestStringsIter(t *testing.T) {
	s := lstr.NewStrings([]string{"", "this ", "is ", " ", " a ", "test"})
	var i liter.Iter[string] = s
	expect := slice.NewIter([]string{"this", "is", "a", "test"})
	expectIdx := slice.NewIter([]int{1, 2, 4, 5})
	liter.ForIdx(i, func(str string, sIdx int) {
		assert.Equal(t, expect.Pop(), str)
		assert.Equal(t, expectIdx.Pop(), sIdx)
	})

	i.(liter.Starter[string]).Start()
	str, _ := i.Cur()
	assert.Equal(t, "this", str)
	assert.Equal(t, 1, i.Idx())
}
