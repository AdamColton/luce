package slice_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	s := slice.New([]int{1, 2, 3})
	assert.Equal(t, slice.Slice[int]{1, 2, 3}, s)
}
