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

func TestConsts(t *testing.T) {
	i := ints.MaxI
	i++
	assert.Equal(t, ints.MinI, i)

	i8 := ints.MaxI8
	i8++
	assert.Equal(t, ints.MinI8, i8)

	i16 := ints.MaxI16
	i16++
	assert.Equal(t, ints.MinI16, i16)

	i32 := ints.MaxI32
	i32++
	assert.Equal(t, ints.MinI32, i32)

	i64 := ints.MaxI64
	i64++
	assert.Equal(t, ints.MinI64, i64)

	u := ints.MaxU
	u++
	assert.Equal(t, uint(0), u)

	u8 := ints.MaxU8
	u8++
	assert.Equal(t, uint8(0), u8)

	u16 := ints.MaxU16
	u16++
	assert.Equal(t, uint16(0), u16)

	u32 := ints.MaxU32
	u32++
	assert.Equal(t, uint32(0), u32)

	u64 := ints.MaxU64
	u64++
	assert.Equal(t, uint64(0), u64)
}

func TestGCD(t *testing.T) {
	expected := 5
	a := 2 * expected
	b := 3 * expected
	assert.Equal(t, expected, ints.GCD(a, b))
}

func TestLCM(t *testing.T) {
	expected := 2 * 3 * 5
	a := 2 * 5
	b := 3 * 5
	assert.Equal(t, expected, ints.LCM(a, b))
}
