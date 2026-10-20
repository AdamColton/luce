package roaring_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lset/roaring"
	"github.com/stretchr/testify/assert"
)

func TestArray(t *testing.T) {
	// Track which days of the year (0-365) have a meeting.
	var meetings roaring.Array
	assert.Equal(t, 0, meetings.Len())
	assert.False(t, meetings.Contains(45))

	// Added out of order; the set keeps them sorted internally.
	meetings.Add(200, 45, 310)
	assert.True(t, meetings.Checksert(45), "already booked")
	assert.False(t, meetings.Checksert(0), "newly booked")
	assert.Equal(t, 4, meetings.Len())

	assert.True(t, meetings.Contains(0))
	assert.True(t, meetings.Contains(45))
	assert.True(t, meetings.Contains(200))
	assert.True(t, meetings.Contains(310))
	assert.False(t, meetings.Contains(100))

	meetings.Remove(45).Remove(100)
	assert.False(t, meetings.Contains(45))
	assert.True(t, meetings.Contains(200))
	assert.Equal(t, 3, meetings.Len())

	// Insert and Delete report whether anything changed.
	assert.True(t, meetings.Insert(150))
	assert.False(t, meetings.Insert(150), "already booked")
	assert.True(t, meetings.Delete(150))
	assert.False(t, meetings.Delete(150), "already cancelled")
}

func TestArrayEdges(t *testing.T) {
	var a roaring.Array
	a.Remove(7) // from an empty set
	a.Add(65535, 0)
	assert.True(t, a.Contains(0))
	assert.True(t, a.Contains(65535))
	assert.False(t, a.Contains(1))
	a.Remove(0).Remove(65535)
	assert.Equal(t, 0, a.Len())
}

func TestArrayMatchesBitmap(t *testing.T) {
	// Array and Bitmap hold the same set, differently; driven through the
	// same changes they always agree.
	var a roaring.Array
	var b roaring.Bitmap
	for i := range 5000 {
		v := uint16(i * 7919) // a prime stride, so values arrive out of order
		assert.Equal(t, b.Checksert(v), a.Checksert(v))
		if i%3 == 0 {
			r := uint16(i * 104729)
			a.Remove(r)
			b.Remove(r)
		}
	}
	assert.Equal(t, b.Len(), a.Len())
	for v := range 1 << 16 {
		assert.Equal(t, b.Contains(uint16(v)), a.Contains(uint16(v)))
	}
}
