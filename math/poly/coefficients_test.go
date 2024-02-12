package poly_test

import (
	"testing"

	"github.com/adamcolton/luce/math/poly"
	"github.com/stretchr/testify/assert"
)

func TestSlice(t *testing.T) {
	s := poly.Slice{1, 2, 3}
	assert.Equal(t, 3, s.Len())
	assert.Equal(t, 2.0, s.AtIdx(1))
	assert.Equal(t, 0.0, s.AtIdx(3), "past the end is 0")

	// the buffers are reused when they have the capacity
	buf := make([]float64, 2, 10)
	b := poly.Buf(5, buf)
	assert.Equal(t, poly.Slice{1}, b, "the constant polynomial 1")
	assert.Equal(t, 10, cap(b))
	assert.Equal(t, 20, cap(poly.Buf(20, buf)))

	assert.Len(t, poly.BufLen(buf, 4), 4)
	assert.Equal(t, 10, cap(poly.BufLen(buf, 4)))
	assert.Len(t, poly.BufLen(nil, 3), 3)
	assert.Len(t, poly.BufEmpty(buf, 4), 0)
	assert.Equal(t, 10, cap(poly.BufEmpty(buf, 4)))

	s, rest := poly.BufSplit(make([]float64, 5), 2)
	assert.Equal(t, 2, cap(s))
	assert.Equal(t, 3, cap(rest))
}
func TestConstants(t *testing.T) {
	var e poly.Coefficients = poly.Empty{}
	assert.Equal(t, 0, e.Len())
	assert.Equal(t, 0.0, e.AtIdx(0))

	var d0 poly.Coefficients = poly.D0(5)
	assert.Equal(t, 1, d0.Len())
	assert.Equal(t, 5.0, d0.AtIdx(0))
	assert.Equal(t, 0.0, d0.AtIdx(1))

	var d1 poly.Coefficients = poly.D1(5)
	assert.Equal(t, 2, d1.Len())
	assert.Equal(t, 5.0, d1.AtIdx(0))
	assert.Equal(t, 1.0, d1.AtIdx(1))
	assert.Equal(t, 0.0, d1.AtIdx(2))
}
