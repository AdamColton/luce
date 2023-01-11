package slice_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

func TestBufferEmpty(t *testing.T) {
	var buf slice.Slice[float64]
	buf = buf.Buffer().Empty(10)
	assert.Equal(t, 10, cap(buf))

	buf = buf.Buffer().Empty(5)
	assert.True(t, cap(buf) >= 5)

	buf = buf.Buffer().Empty(12)
	assert.True(t, cap(buf) >= 12)
}

type mockLener int

func (m mockLener) Len() int {
	return int(m)
}

func TestBufferLener(t *testing.T) {
	var buf slice.Buffer[float64]
	s := buf.Lener("not a Lener")
	assert.Equal(t, 0, cap(s))

	l := mockLener(10)
	s = buf.Lener(l)
	assert.Equal(t, 10, cap(s))
}

func TestBufferSlice(t *testing.T) {
	buf := (slice.Buffer[float64]{3, 1, 4})[:0]
	s := buf.Slice(2)

	assert.Equal(t, slice.Slice[float64]{3, 1}, s)
	assert.Equal(t, slice.Slice[float64]{0, 0, 0, 0}, buf.Slice(4))
}

func TestBufferZeros(t *testing.T) {
	buf := (slice.Buffer[float64]{3, 1, 4, 1, 5})[:3]
	s := buf.Zeros(5)
	assert.Equal(t, slice.Slice[float64]{0, 0, 0, 0, 0}, s)
	s = buf.Zeros(6)
	assert.Equal(t, slice.Slice[float64]{0, 0, 0, 0, 0, 0}, s)
}

func TestReduceCapacity(t *testing.T) {
	buf := slice.Buffer[float64]{1, 2, 3, 4, 5}
	sub := buf[:3]
	assert.Equal(t, 5, cap(sub))
	s := buf.ReduceCapacity(3)
	assert.Equal(t, 3, cap(s))
	assert.Equal(t, 5, cap(buf))
	s = buf.ReduceCapacity(10)
	assert.Equal(t, 5, cap(s))
}

func TestSplit(t *testing.T) {
	buf := slice.Buffer[float64]{}.Empty(15).Buffer()
	assert.Equal(t, 15, cap(buf))

	a, b := buf.Split(5)

	assert.Equal(t, 5, cap(a))
	assert.Equal(t, 10, cap(b))

	c, d := b.Split(12)
	assert.Equal(t, 12, cap(c))
	assert.Equal(t, 10, cap(d))
}
