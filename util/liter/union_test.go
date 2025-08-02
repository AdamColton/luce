package liter_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/util/liter"
	"github.com/stretchr/testify/assert"
)

func TestUnion(t *testing.T) {
	u := liter.NewUnion(
		slice.Slice[int]{3, 1, 4, 1, 5, 9, 2, 6, 5, 3}.Iter(),
		list.NewGenerator(10, func(i int) int { return i }).Iter(),
	)
	assert.False(t, u.Done())
	assert.Equal(t, 0, u.Idx())
	got := slice.FromIter(u, nil)
	assert.True(t, u.Done())
	assert.Equal(t, 10, u.Idx())
	expected := slice.Slice[[]int]{
		{3, 0},
		{1, 1},
		{4, 2},
		{1, 3},
		{5, 4},
		{9, 5},
		{2, 6},
		{6, 7},
		{5, 8},
		{3, 9},
	}
	assert.Equal(t, expected, got)

	u = liter.NewUnion[int]()
	assert.Equal(t, 0, u.Idx())
	assert.True(t, u.Done())
}
