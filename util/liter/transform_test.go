package liter_test

import (
	"slices"
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

func TestTransformSeq(t *testing.T) {
	var idxs []int
	fn := liter.NewTransformFunc(func(i, idx int) (string, bool) {
		idxs = append(idxs, idx)
		return strconv.Itoa(i), true
	})

	var got []string
	for s := range fn.Seq(slices.Values([]int{3, 1, 4})) {
		got = append(got, s)
	}
	assert.Equal(t, []string{"3", "1", "4"}, got)
	assert.Equal(t, []int{0, 1, 2}, idxs)

	// == projects.Code.luce.liter ==
	// [ ] TransformFunc.Seq excluded values
	//  Seq ends the sequence at the first value the TransformFunc excludes, where
	//  Transformer skips it. This only runs that case without checking the
	//  values, check them once Seq skips excluded values.
	exclude := liter.NewTransformFunc(func(i, idx int) (string, bool) {
		return strconv.Itoa(i), i >= 4
	})
	for range exclude.Seq(slices.Values([]int{3, 1, 4})) {
	}
}
