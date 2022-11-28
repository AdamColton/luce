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
