package list_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

func TestPairs(t *testing.T) {
	l := list.NewGenerator(5, func(i int) int {
		return i
	})

	p := list.NewPairs(l, true)
	assert.Equal(t, l.Len(), p.Len())
	p.Loop = false
	assert.Equal(t, l.Len()-1, p.Len())

	w := p.Wrap()
	assert.Equal(t, list.Wrapper[[2]int]{p}, w)

	p.Loop = true
	expected := [][2]int{
		{0, 1},
		{1, 2},
		{2, 3},
		{3, 4},
		{4, 0},
	}
	assert.Equal(t, expected, w.Slice(nil))
}

func TestPairsShorthand(t *testing.T) {
	a := slice.New([]int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3})

	p := list.NewPairs(a, true)
	expected := [][2]int{
		{3, 1},
		{1, 4},
		{4, 1},
		{1, 5},
		{5, 9},
		{9, 2},
		{2, 6},
		{6, 5},
		{5, 3},
		{3, 3},
	}

	var got [][2]int
	p.For(func(v [2]int) {
		got = append(got, v)
	})
	assert.Equal(t, expected, got)

	got = got[:0]
	p.Each(func(idx int, v [2]int, done *bool) {
		assert.Equal(t, len(got), idx)
		got = append(got, v)
	})
	assert.Equal(t, expected, got)
}
