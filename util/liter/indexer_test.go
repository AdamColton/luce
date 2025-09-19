package liter_test

import (
	"testing"

	"github.com/adamcolton/luce/util/liter"
	"github.com/stretchr/testify/assert"
)

func fibNexter(max int) func() (int, bool) {
	cur, prev := 1, 0
	return func() (int, bool) {
		cur, prev = cur+prev, cur
		return cur, cur > max
	}
}

func TestIndexer(t *testing.T) {
	n := liter.NewNextFunc(fibNexter(100))
	var got []int
	i := n.Indexer()
	assert.False(t, i.Done())
	assert.Equal(t, 0, i.Idx())
	liter.For(i, func(i int) {
		got = append(got, i)
	})
	expected := []int{1, 2, 3, 5, 8, 13, 21, 34, 55, 89}
	assert.Equal(t, expected, got)
	assert.True(t, i.Done())
	assert.Equal(t, len(expected), i.Idx())
}
