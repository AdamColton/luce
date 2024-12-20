package ints_test

import (
	"strconv"
	"testing"

	"github.com/adamcolton/luce/math/ints"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/constraints"
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

func TestConversions(t *testing.T) {
	testConversions(t, int(5))
	testConversions(t, int8(5))
	testConversions(t, int16(5))
	testConversions(t, int32(5))
	testConversions(t, int64(5))
	testConversions(t, uint(5))
	testConversions(t, uint8(5))
	testConversions(t, uint16(5))
	testConversions(t, uint32(5))
	testConversions(t, uint64(5))
}

func testConversions[I constraints.Integer](t *testing.T, five I) {
	assert.Equal(t, int(5), ints.Int(five))
	assert.Equal(t, int8(5), ints.Int8(five))
	assert.Equal(t, int16(5), ints.Int16(five))
	assert.Equal(t, int32(5), ints.Int32(five))
	assert.Equal(t, int64(5), ints.Int64(five))
	assert.Equal(t, uint(5), ints.Uint(five))
	assert.Equal(t, uint8(5), ints.Uint8(five))
	assert.Equal(t, uint16(5), ints.Uint16(five))
	assert.Equal(t, uint32(5), ints.Uint32(five))
	assert.Equal(t, uint64(5), ints.Uint64(five))
	assert.Equal(t, float32(5), ints.Float32(five))
	assert.Equal(t, float64(5), ints.Float64(five))
}

func TestReduce(t *testing.T) {
	got := ints.Reduce(ints.GCD, []int{6, 8, 10, 12})
	assert.Equal(t, 2, got)

	got = ints.Reduce(ints.GCD, []int{})
	assert.Equal(t, 0, got)
}

func TestLCMN(t *testing.T) {
	s := []int{
		2 * 5 * 7,
		2 * 3 * 7,
		2 * 3 * 5,
		3 * 5 * 7,
	}
	assert.Equal(t, 2*3*5*7, ints.LCMN(s...))
}

func TestProd(t *testing.T) {
	assert.Equal(t, 1, ints.Prod[int]())
	assert.Equal(t, 24, ints.Prod(1, 2, 3, 4))
}

func TestCompoundZero(t *testing.T) {
	assert.Equal(t, 0, ints.Reduce(ints.SumFn, []int{}))
	assert.Equal(t, 1, ints.Reduce(ints.SumFn, []int{1}))
	assert.Equal(t, 3, ints.Reduce(ints.SumFn, []int{1, 2}))
	assert.Equal(t, 6, ints.Reduce(ints.SumFn, []int{1, 2, 3}))
}

func TestRange(t *testing.T) {
	assert.Equal(t, 5, ints.Range(0, 5, 10))
	assert.Equal(t, 0, ints.Range(0, -1, 10))
	assert.Equal(t, 10, ints.Range(0, 11, 10))
}

func TestIdx(t *testing.T) {
	tt := map[int]struct {
		expected int
		ok       bool
	}{
		0: {
			expected: 0,
			ok:       true,
		},
		4: {
			expected: 4,
			ok:       true,
		},
		5: {
			expected: 5,
			ok:       false,
		},
		-1: {
			expected: 4,
			ok:       true,
		},
		-5: {
			expected: 0,
			ok:       true,
		},
		-6: {
			expected: -1,
			ok:       false,
		},
	}

	for n, tc := range tt {
		t.Run(strconv.Itoa(n), func(t *testing.T) {
			idx, ok := ints.Idx(n, 5)
			assert.Equal(t, tc.expected, idx)
			assert.Equal(t, tc.ok, ok)
		})
	}
}
