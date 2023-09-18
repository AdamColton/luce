package filter_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

func TestFilterOps(t *testing.T) {
	gt4 := filter.Filter[int](func(i int) bool {
		return i > 4
	})
	lt8 := filter.Filter[int](func(i int) bool {
		return i < 8
	})

	gt4AndLt8 := gt4.And(lt8)
	assert.False(t, gt4AndLt8(3))
	assert.True(t, gt4AndLt8(5))
	assert.False(t, gt4AndLt8(10))

	lte4OrGte8 := gt4.Not().Or(lt8.Not())
	assert.True(t, lte4OrGte8(3))
	assert.False(t, lte4OrGte8(5))
	assert.True(t, lte4OrGte8(10))
}

func TestSlice(t *testing.T) {
	gt4 := filter.Filter[int](func(i int) bool {
		return i > 4
	})
	s := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5}
	got := gt4.Slice(s)
	expected := slice.Slice[int]{5, 9, 6, 5, 5}
	assert.Equal(t, expected, got)

	gt10 := filter.Filter[int](func(i int) bool {
		return i > 10
	})
	got = gt10.Slice(s)
	assert.Nil(t, got)
}

func TestChan(t *testing.T) {
	gt4 := filter.Filter[int](func(i int) bool {
		return i > 4
	})
	ch := make(chan int)
	go func() {
		for _, i := range []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5} {
			ch <- i
		}
		close(ch)
	}()

	to := timeout.After(5, func() {
		expected := []int{5, 9, 6, 5, 5}
		get := gt4.Chan(ch, 0)
		idx := 0
		for got := range get {
			assert.Equal(t, expected[idx], got)
			idx++
		}
	})
	assert.NoError(t, to)
}
