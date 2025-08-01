package list_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/list"
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
