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
	c := 0
	testFunc := func(idx int, str string, done *bool) {
		c++
		assert.Equal(t, strconv.Itoa(s.Slice[idx]), str)
	}
	tf.Each(testFunc)
	assert.Len(t, s.Slice, c)
	assert.True(t, tf.Done())

	c = 0
	fn.Factory(fac).Each(testFunc)
	assert.Len(t, s.Slice, c)
	assert.True(t, tf.Done())

	fn2 := liter.NewTransformFunc(func(i, idx int) (string, bool) {
		if i < 4 {
			return "", false
		}
		return strconv.Itoa(i), true
	})
	expected := []string{"4", "5", "9"}
	s.idx = 0
	c = 0
	fn2.New(s).Each(func(idx int, str string, done *bool) {
		c++
		assert.Equal(t, expected[idx], str)
	})
	assert.Len(t, expected, c)
}
