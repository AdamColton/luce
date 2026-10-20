package roaring_test

import (
	"slices"
	"testing"

	"github.com/adamcolton/luce/ds/lset/roaring"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/util/liter"
	"github.com/adamcolton/luce/util/upgrade"
	"github.com/stretchr/testify/assert"
)

func TestSetSlice(t *testing.T) {
	// Values come back in ascending order, whatever order they went in.
	blocked := roaring.New(ip(203, 0, 113, 9), ip(10, 0, 0, 2), ip(192, 168, 1, 7), ip(10, 0, 0, 1))
	want := []uint32{ip(10, 0, 0, 1), ip(10, 0, 0, 2), ip(192, 168, 1, 7), ip(203, 0, 113, 9)}
	assert.Equal(t, want, []uint32(blocked.Slice(nil)))

	// A buffer with enough capacity is reused.
	buf := make([]uint32, 0, 10)
	got := blocked.Slice(buf)
	assert.Equal(t, want, []uint32(got))
	assert.Equal(t, &buf[:1][0], &got[0])

	var nilSet *roaring.Set
	assert.Nil(t, nilSet.Slice(nil))
}

func TestSetIter(t *testing.T) {
	// Walk a blocklist in order. 172.16.0.0/16 is fully blocked, so its
	// container is a Bitmap; the other networks hold a few addresses each.
	blocked := roaring.New(ip(10, 0, 0, 1), ip(192, 168, 1, 7))
	for c := range 256 {
		for d := range 256 {
			blocked.Add(ip(172, 16, byte(c), byte(d)))
		}
	}

	var got []uint32
	for v := range blocked.Iter().Seq {
		got = append(got, v)
	}
	assert.Equal(t, []uint32(blocked.Slice(nil)), got)
	assert.Len(t, got, 2+65536)
	assert.True(t, slices.IsSorted(got))
	assert.Equal(t, ip(10, 0, 0, 1), got[0])
	assert.Equal(t, ip(172, 16, 0, 0), got[1])
	assert.Equal(t, ip(172, 16, 255, 255), got[65536])
	assert.Equal(t, ip(192, 168, 1, 7), got[65537])
}

func TestSetIterMethods(t *testing.T) {
	// Some scattered bits in one Bitmap, between two Array containers.
	s := roaring.New(5, 1<<16|3)
	for v := range uint32(5000) {
		s.Add(2<<16 | v*13)
	}
	s.Add(3<<16 | 0xffff)
	want := s.Slice(nil)

	it := s.Iter()
	v, done := it.Cur()
	assert.False(t, done)
	assert.Equal(t, uint32(5), v)
	assert.Equal(t, 0, it.Idx())
	for idx := 1; idx < len(want); idx++ {
		v, done = it.Next()
		assert.False(t, done)
		assert.Equal(t, want[idx], v)
		assert.Equal(t, idx, it.Idx())
	}
	_, done = it.Next()
	assert.True(t, done)
	assert.True(t, it.Done())
	_, done = it.Next()
	assert.True(t, done, "stays done")

	// slice.FromIter takes the Slice fast path; the whole set, whatever the
	// iterator's position.
	assert.Equal(t, want, slice.FromIter[uint32](it, nil))

	starter, ok := upgrade.To[liter.Starter[uint32]](it)
	assert.True(t, ok)
	v, done = starter.Start()
	assert.False(t, done)
	assert.Equal(t, uint32(5), v)
	assert.Equal(t, 0, it.Idx())
}

func TestSetIterEmpty(t *testing.T) {
	var nilSet *roaring.Set
	assert.True(t, nilSet.Iter().Done())

	s := roaring.New(7)
	s.Remove(7)
	_, done := s.Iter().Cur()
	assert.True(t, done)
}
