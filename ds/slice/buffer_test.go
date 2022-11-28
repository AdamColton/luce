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
