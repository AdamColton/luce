package lstr

import (
	"testing"

	"github.com/adamcolton/luce/math/ints"
	"github.com/stretchr/testify/assert"
)

// A string long enough to overflow an int cannot be created, so the overflow
// check is tested directly.
func TestAddLen(t *testing.T) {
	assert.Equal(t, 5, addLen(2, 3))
	assert.Equal(t, ints.MaxI, addLen(ints.MaxI-1, 1))
	assert.Panics(t, func() { addLen(ints.MaxI, 1) })
	assert.Panics(t, func() { addLen(ints.MaxI-1, 2) })
}
