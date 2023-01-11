package slice_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

func TestBufferEmpty(t *testing.T) {
	var buf []float64
	buf = slice.BufferEmpty(10, buf)
	assert.Equal(t, 10, cap(buf))

	buf = slice.BufferEmpty(5, buf)
	assert.True(t, cap(buf) >= 5)

	buf = slice.BufferEmpty(12, buf)
	assert.True(t, cap(buf) >= 12)
}

type mockLener int

func (m mockLener) Len() int {
	return int(m)
}

func TestBufferLener(t *testing.T) {
	var buf []float64
	buf = slice.BufferLener("not a Lener", buf)
	assert.Equal(t, 0, cap(buf))

	l := mockLener(10)
	buf = slice.BufferLener(l, buf)
	assert.Equal(t, 10, cap(buf))
}

func TestBufferSlice(t *testing.T) {
	buf := ([]float64{3, 1, 4})[:0]
	buf = slice.BufferSlice(2, buf)

	assert.Equal(t, []float64{3, 1}, buf)
	assert.Equal(t, []float64{0, 0, 0, 0}, slice.BufferSlice(4, buf))
}

func TestBufferZeros(t *testing.T) {
	buf := ([]float64{3, 1, 4, 1, 5})[:3]
	buf = slice.BufferZeros(5, buf)
	assert.Equal(t, []float64{0, 0, 0, 0, 0}, buf)
	buf = slice.BufferZeros(6, buf)
	assert.Equal(t, []float64{0, 0, 0, 0, 0, 0}, buf)
}

func TestReduceCapacity(t *testing.T) {
	buf := []float64{1, 2, 3, 4, 5}
	sub := buf[:3]
	assert.Equal(t, 5, cap(sub))
	sub = slice.ReduceCapacity(3, sub)
	assert.Equal(t, 3, cap(sub))
	assert.Equal(t, 5, cap(buf))
}

func TestSplit(t *testing.T) {
	buf := slice.BufferEmpty[float64](15, nil)
	a, b := slice.BufferSplit(5, buf)

	assert.Equal(t, 5, cap(a))
	assert.Equal(t, 10, cap(b))

	c, d := slice.BufferSplit(12, b)
	assert.Equal(t, 12, cap(c))
	assert.Equal(t, 10, cap(d))
}
