package filter_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/liter"
	"github.com/stretchr/testify/assert"
)

func TestIter(t *testing.T) {
	var iterOut liter.Iter[int]
	s := slice.Slice[int]{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5, 8}

	iterOut = filter.GTE(5).Iter(s.Iter())

	expected := []int{5, 9, 6, 5, 5, 8}
	checkFn := func(i, idx int) {
		assert.Equal(t, expected[idx], i, idx)
	}
	c := liter.Wrap(iterOut).ForIdx(checkFn)
	assert.Len(t, expected, c)
	assert.True(t, iterOut.Done())

	iterOut = filter.LT(5).Iter(s.Iter())
	expected = []int{3, 1, 4, 1, 2, 3}
	c = liter.Wrap(iterOut).ForIdx(checkFn)
	assert.Len(t, expected, c)
	assert.True(t, iterOut.Done())
}
