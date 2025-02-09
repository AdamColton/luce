package liter_test

import (
	"strconv"
	"testing"

	"github.com/adamcolton/luce/util/liter"
	"github.com/stretchr/testify/assert"
)

func TestTransform(t *testing.T) {
	s := &sliceIter[int]{
		Slice: []int{3, 1, 4, 1, 5, 9},
	}
	fac := func() (liter.Iter[int], int, bool) {
		s.idx = 0
		return s, s.Slice[0], false
	}

	fn := liter.ForAll(strconv.Itoa)
	tf := fn.New(s)
	testFunc := func(str string, idx int) {
		assert.Equal(t, strconv.Itoa(s.Slice[idx]), str)
	}
	i := tf.ForIdx(testFunc)
	assert.True(t, i == len(s.Slice)-1)
	assert.True(t, tf.Done())

	i = fn.Factory(fac).ForIdx(testFunc)
	assert.True(t, i == len(s.Slice)-1)
	assert.True(t, tf.Done())

	fn = func(i, idx int) (string, bool) {
		if i < 4 {
			return "", false
		}
		return strconv.Itoa(i), true
	}
	expected := []string{"4", "5", "9"}
	s.idx = 0
	i = fn.New(s).ForIdx(func(str string, idx int) {
		assert.Equal(t, expected[idx], str)
	})
	assert.True(t, i == 2)
}
