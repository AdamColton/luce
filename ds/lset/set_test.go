package lset_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lset"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	lt := slice.LT[int]()
	s := lset.New[int]()
	s.Add(3, 1, 4, 1, 5, 9)

	assert.True(t, s.Contains(1))
	assert.False(t, s.Contains(2))

	s.Remove(1)
	assert.False(t, s.Contains(1))

	expect := func(expected ...int) {
		assert.Equal(t, expected, lt.Sort(s.Slice(nil)))
	}
	expect(3, 4, 5, 9)

	assert.Equal(t, s.Len(), 4)

	s2 := lset.New(6, 7)
	s.AddAll(s2)
	expect(3, 4, 5, 6, 7, 9)

	s2 = s.Copy()
	assert.Equal(t, s, s2)
}

func TestMulti(t *testing.T) {
	s1 := lset.New(1, 2, 3, 4)
	s2 := lset.New(4, 5)
	s3 := lset.New(5, 6, 7)
	m := lset.NewMulti(s1, s2, s3)

	m.Sort()
	assert.Equal(t, lset.Multi[int]{s2, s3, s1}, m)

	assert.True(t, m.Contains(3))
	assert.False(t, m.Contains(10))
}
