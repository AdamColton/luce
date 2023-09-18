package filter_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/channel"
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
	s := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5}
	snd := channel.Slice(s, nil)
	// since the channel is populated, close can be called immediatly and it will
	// close when data is drained.
	close(snd)
	p, _, rcv := channel.NewPipe(snd, nil)

	expected := []int{5, 9, 6, 5, 5}
	idx := 0
	to := timeout.After(5, func() {
		gt4.Chan(p)
		for got := range rcv {
			assert.Equal(t, expected[idx], got)
			idx++
		}
	})
	assert.Equal(t, len(expected), idx)
	assert.NoError(t, to)
}
