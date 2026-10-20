package roaring_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lset/roaring"
	"github.com/stretchr/testify/assert"
)

func TestBitmap(t *testing.T) {
	// Track which TCP ports are open.
	var open roaring.Bitmap
	assert.Equal(t, 0, open.Len())
	assert.False(t, open.Contains(80))

	open.Add(22, 80, 443)
	assert.True(t, open.Checksert(80), "already open")
	assert.False(t, open.Checksert(8080), "newly opened")
	assert.Equal(t, 4, open.Len())

	assert.True(t, open.Contains(22))
	assert.True(t, open.Contains(443))
	assert.False(t, open.Contains(21))

	open.Remove(22).Remove(8080).Remove(25)
	assert.False(t, open.Contains(22))
	assert.Equal(t, 2, open.Len())

	// Insert and Delete report whether anything changed.
	assert.True(t, open.Insert(3306))
	assert.False(t, open.Insert(3306), "already open")
	assert.True(t, open.Delete(3306))
	assert.False(t, open.Delete(3306), "already closed")
}

func TestBitmapEdges(t *testing.T) {
	// The first and last value, and both ends of a word.
	vals := []uint16{0, 63, 64, 65535}
	var b roaring.Bitmap
	b.Add(vals...)
	assert.Equal(t, len(vals), b.Len())
	for _, v := range vals {
		assert.True(t, b.Contains(v))
	}
	assert.False(t, b.Contains(1))
	assert.False(t, b.Contains(62))
	assert.False(t, b.Contains(65534))
}

func TestBitmapFull(t *testing.T) {
	var b roaring.Bitmap
	for v := range 1 << 16 {
		b.Add(uint16(v))
	}
	assert.Equal(t, 1<<16, b.Len())
	for v := range 1 << 16 {
		b.Remove(uint16(v))
	}
	assert.Equal(t, 0, b.Len())
}
