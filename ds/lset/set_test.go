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
		assert.Equal(t, expected, lt.Sort(s.Slice()))
	}
	expect(3, 4, 5, 9)

	assert.Equal(t, s.Len(), 4)

	s2 := lset.New(6, 7)
	s.AddAll(s2)
	expect(3, 4, 5, 6, 7, 9)

	s2 = s.Copy()
	assert.Equal(t, s, s2)

	got := make([]int, 0, s.Len())
	s.Each(func(i int, done *bool) {
		got = append(got, i)
	})
	assert.Equal(t, []int{3, 4, 5, 6, 7, 9}, slice.LT[int]().Sort(got))

	got = got[:0]
	s.Each(func(i int, done *bool) {
		got = append(got, i)
		*done = len(got) >= 3
	})
	assert.Len(t, got, 3)
}

func TestEachNoNilPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error(r)
		}
	}()
	var s *lset.Set[string]
	s.Each(func(str string, done *bool) {
		t.Error("this should not be reached")
	})
}

func TestMulti(t *testing.T) {
	s1 := lset.New(1, 2, 3, 4, 8)
	s2 := lset.New(1, 4, 5)
	s3 := lset.New(1, 5, 6, 7)
	m := lset.NewMulti(s1, s2, s3)

	m.Sort()
	assert.Equal(t, lset.Multi[int]{s2, s3, s1}, m)

	assert.True(t, m.Contains(3))
	assert.False(t, m.Contains(10))

	assert.True(t, m.AllContain(1))
	assert.False(t, m.AllContain(5))
}

func TestMultiIntersection(t *testing.T) {
	s1 := lset.New(1, 2, 3, 4, 5)
	s2 := lset.New(1, 2, 6, 7)
	s3 := lset.New(1, 2, 8, 9)
	m := lset.Multi[int]{s1, s2, s3}

	m.Sort()
	i := m.Intersection()
	got := slice.LT[int]().Sort(i.Slice())
	assert.Equal(t, []int{1, 2}, got)

	m = lset.Multi[int]{}
	assert.Nil(t, m.Intersection())

	m = lset.Multi[int]{s1}
	assert.Equal(t, s1, m.Intersection())
}
