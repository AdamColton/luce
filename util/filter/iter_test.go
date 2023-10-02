package filter_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/liter"
	"github.com/stretchr/testify/assert"
)

func TestIter(t *testing.T) {
	var iterIn, iterOut liter.Iter[int]
	iterIn = (slice.Slice[int]{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5, 8}).Iter()

	iterOut = filter.GTE(5).Iter(iterIn)

	expected := []int{5, 9, 6, 5, 5, 8}
	c := liter.Wrap(iterOut).ForIdx(func(i, idx int) {
		assert.Equal(t, expected[idx], i, idx)
	})
	assert.Len(t, expected, c)
	assert.True(t, iterOut.Done())
}
