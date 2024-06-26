package liter_test

import (
	"testing"

	"github.com/adamcolton/luce/util/liter"
	"github.com/stretchr/testify/assert"
)

func TestReduce(t *testing.T) {
	max := liter.Reducer[int, int](func(aggregate, element, idx int) int {
		if idx == 0 || element > aggregate {
			return element
		}
		return aggregate
	})
	s := []int{3, 1, 4, 1, 5, 9}

	si := &sliceIter[int]{
		Slice: s,
	}
	assert.Equal(t, 9, max.Iter(0, si))

	var sf liter.Factory[int] = sliceFactory(s)
	assert.Equal(t, 9, max.Factory(100, sf))
}

func TestAppend(t *testing.T) {
	s := &sliceIter[int]{
		Slice: []int{3, 1, 4, 1, 5, 9},
	}
	got := liter.Appender[int]().Iter(nil, s)
	assert.Equal(t, s.Slice, got)
}

func TestMinMax(t *testing.T) {
	s := &sliceIter[string]{
		Slice: []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine"},
	}
	lnFn := func(str string) int { return len(str) }

	got := liter.Min(lnFn).Iter(100, s)
	assert.Equal(t, 3, got)

	s.idx = 0
	got = liter.Max(lnFn).Iter(0, s)
	assert.Equal(t, 5, got)
}
