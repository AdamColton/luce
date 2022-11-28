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
