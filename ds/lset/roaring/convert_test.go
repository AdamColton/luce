package roaring_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lset/roaring"
	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	// A sparse set starts as an Array; once it grows dense it is cheaper as
	// a Bitmap, and it can go back again if it shrinks.
	var a roaring.Array
	a.Add(1000, 3, 65535, 0, 64, 63)

	b := a.Bitmap(nil)
	assert.Equal(t, a.Len(), b.Len())
	for _, v := range []uint16{0, 3, 63, 64, 1000, 65535} {
		assert.True(t, b.Contains(v))
	}
	assert.False(t, b.Contains(1))

	back := b.Array(nil)
	assert.Equal(t, a, back)
}

func TestConvertBuffers(t *testing.T) {
	var a roaring.Array
	a.Add(10, 20)

	// A Bitmap buffer is cleared before use, so earlier contents don't leak.
	buf := &roaring.Bitmap{}
	buf.Add(500)
	b := a.Bitmap(buf)
	assert.Same(t, buf, b)
	assert.False(t, b.Contains(500))
	assert.Equal(t, 2, b.Len())

	// A slice buffer with enough capacity is reused.
	values := make([]uint16, 0, 8)
	back := b.Array(values)
	assert.Equal(t, a, back)
	back.Add(5)
	assert.Equal(t, uint16(5), values[:1][0], "written into buf's memory")
}

func TestConvertEdges(t *testing.T) {
	var empty roaring.Array
	assert.Equal(t, 0, empty.Bitmap(nil).Len())
	var emptyBitmap roaring.Bitmap
	emptyArray := emptyBitmap.Array(nil)
	assert.Equal(t, 0, emptyArray.Len())

	var full roaring.Bitmap
	for v := range 1 << 16 {
		full.Add(uint16(v))
	}
	a := full.Array(nil)
	assert.Equal(t, 1<<16, a.Len())
	assert.Equal(t, full, *a.Bitmap(nil))
}
