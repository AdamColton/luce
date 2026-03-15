package slice_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

func TestReverseReturnsSlice(t *testing.T) {
	s := slice.Slice[int]{1, 2, 3}
	assert.Equal(t, slice.Slice[int]{3, 2, 1}, s.Reverse())
	assert.Equal(t, slice.Slice[int]{3, 2, 1}, s, "the slice itself is reversed")
}
