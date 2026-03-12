package liter_test

import (
	"testing"

	"github.com/adamcolton/luce/util/liter"
	"github.com/stretchr/testify/assert"
)

func TestBest(t *testing.T) {
	s := &sliceIter[[2]int]{
		Slice: [][2]int{
			{3, 1},
			{4, 1},
			{5, 9},
			{2, 6},
			{5, 3},
			{5, 8},
			{9, 7},
			{3, 2},
		},
	}

	fn := func(i [2]int) int {
		return i[0] * i[1]
	}

	p, c, idx := liter.Greatest(s, fn)
	assert.Equal(t, [2]int{9, 7}, p)
	assert.Equal(t, 63, c)
	assert.Equal(t, 6, idx)

	s.idx = 0
	p, c, idx = liter.Least(s, fn)
	assert.Equal(t, [2]int{3, 1}, p)
	assert.Equal(t, 3, c)
	assert.Equal(t, 0, idx)
}
