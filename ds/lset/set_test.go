package lset_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lset"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	s := lset.New[int]()
	s.Add(3, 1, 4, 1, 5, 9)

	assert.True(t, s.Contains(1))
	assert.False(t, s.Contains(2))

	s.Remove(1)
	assert.False(t, s.Contains(1))

	assert.Equal(t, []int{3, 4, 5, 9}, slice.LT[int]().Sort(s.Slice()))
	assert.Equal(t, s.Len(), 4)

	s2 := lset.New(6, 7)
	s.AddAll(s2)
	assert.Equal(t, []int{3, 4, 5, 6, 7, 9}, slice.LT[int]().Sort(s.Slice()))

	s2 = s.Copy()
	assert.Equal(t, s, s2)
}
