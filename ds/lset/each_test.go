package lset_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lset"
	"github.com/stretchr/testify/assert"
)

func TestEachDone(t *testing.T) {
	s := lset.New(1, 2, 3, 4)
	calls := 0
	s.Each(func(i int, done *bool) {
		calls++
		*done = true
	})
	assert.Equal(t, 1, calls)

	sum := 0
	s.All(func(i int) { sum += i })
	assert.Equal(t, 10, sum)
}

func TestSortedEach(t *testing.T) {
	s := lset.New(3, 1, 4, 2)
	less := func(i, j int) bool { return i < j }

	var got []int
	s.SortedEach(less, nil, func(i int, done *bool) {
		got = append(got, i)
	})
	assert.Equal(t, []int{1, 2, 3, 4}, got)

	got = nil
	s.SortedEach(less, nil, func(i int, done *bool) {
		got = append(got, i)
		*done = i == 2
	})
	assert.Equal(t, []int{1, 2}, got)

	var nilSet *lset.Set[int]
	nilSet.SortedEach(less, nil, func(i int, done *bool) { t.Fatal("called") })
}
