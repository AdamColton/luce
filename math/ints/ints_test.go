package ints_test

import (
	"testing"

	"github.com/adamcolton/luce/math/ints"
	"github.com/stretchr/testify/assert"
)

func TestDiv(t *testing.T) {
	// 5 items, 2 to a page, fill 3 pages and 2 of them are full.
	assert.Equal(t, 3, ints.DivUp(5, 2))
	assert.Equal(t, 2, ints.DivDown(5, 2))

	tt := []struct {
		a, b, down, up, in, out int
	}{
		{a: 7, b: 2, down: 3, up: 4, in: 3, out: 4},
		{a: -7, b: 2, down: -4, up: -3, in: -3, out: -4},
		{a: 7, b: -2, down: -4, up: -3, in: -3, out: -4},
		{a: -7, b: -2, down: 3, up: 4, in: 3, out: 4},
		{a: 6, b: 2, down: 3, up: 3, in: 3, out: 3},
		{a: -6, b: 2, down: -3, up: -3, in: -3, out: -3},
		{a: 0, b: -2, down: 0, up: 0, in: 0, out: 0},
	}
	for _, tc := range tt {
		assert.Equal(t, tc.down, ints.DivDown(tc.a, tc.b), "DivDown(%d, %d)", tc.a, tc.b)
		assert.Equal(t, tc.up, ints.DivUp(tc.a, tc.b), "DivUp(%d, %d)", tc.a, tc.b)
		assert.Equal(t, tc.in, ints.DivIn(tc.a, tc.b), "DivIn(%d, %d)", tc.a, tc.b)
		assert.Equal(t, tc.out, ints.DivOut(tc.a, tc.b), "DivOut(%d, %d)", tc.a, tc.b)
		assert.Equal(t, tc.down, ints.Div(tc.a, tc.b, ints.Down), "Div(%d, %d, Down)", tc.a, tc.b)
	}

	assert.Equal(t, uint8(3), ints.DivDown[uint8](7, 2))
	assert.Equal(t, uint8(4), ints.DivUp[uint8](7, 2))
	assert.Equal(t, uint8(3), ints.DivIn[uint8](7, 2))
	assert.Equal(t, uint8(4), ints.DivOut[uint8](7, 2))
	assert.Equal(t, int8(-4), ints.DivOut[int8](-7, 2))

	assert.Equal(t, -3, ints.Div(-7, 2, ints.Dir(9)), "an unknown Dir rounds like In")
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
