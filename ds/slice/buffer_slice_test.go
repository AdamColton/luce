package slice_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

func TestSliceBuffer(t *testing.T) {
	s := slice.Slice[int]{1, 2, 3}
	buf := s.Buffer()
	assert.Equal(t, slice.Buffer[int]{1, 2, 3}, buf)

	// The Buffer shares memory with the Slice.
	buf[0] = 10
	assert.Equal(t, 10, s[0])
}
