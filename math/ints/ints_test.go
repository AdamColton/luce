package ints_test

import (
	"testing"

	"github.com/adamcolton/luce/math/ints"
	"github.com/stretchr/testify/assert"
)

func TestDiv(t *testing.T) {
	assert.Equal(t, 3, ints.DivUp(5, 2))
	assert.Equal(t, 2, ints.DivDown(5, 2))
}

func TestMod(t *testing.T) {
	assert.Equal(t, 3, ints.Mod(8, 5))
	assert.Equal(t, 4, ints.Mod(-1, 5))
	assert.Equal(t, 3, ints.Mod(-2, 5))
	assert.Equal(t, 0, ints.Mod(-5, 5))

	assert.Equal(t, -4, ints.Mod(1, -5))
	assert.Equal(t, -3, ints.Mod(2, -5))
	assert.Equal(t, 0, ints.Mod(5, -5))

	assert.Equal(t, -1, ints.Mod(-1, -5))
	assert.Equal(t, -2, ints.Mod(-2, -5))
	assert.Equal(t, 0, ints.Mod(-5, -5))

	assert.Equal(t, -2, (-2 % -5))

}
