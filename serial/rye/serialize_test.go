package rye_test

import (
	"bytes"
	"crypto/rand"
	"math"
	"reflect"
	"sort"
	"strconv"
	"testing"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/serial/rye"
	"github.com/stretchr/testify/assert"
)

func TestSerializeUint16(t *testing.T) {
	tt := []struct {
		expected []byte
		x        uint16
	}{
		{
			expected: []byte{1, 0},
			x:        1,
		}, {
			expected: []byte{255, 0},
			x:        255,
		}, {
			expected: []byte{0, 1},
			x:        256,
		}, {
			expected: []byte{255, 255},
			x:        65535,
		},
	}

	for _, tc := range tt {
		t.Run(strconv.Itoa(int(tc.x)), func(t *testing.T) {
			b := make([]byte, 2)
			rye.Serialize.Uint16(b, tc.x)
			assert.Equal(t, tc.expected, b)

			x := rye.Deserialize.Uint16(tc.expected)
			assert.Equal(t, tc.x, x)
		})
	}
}

func TestSerializeUint32(t *testing.T) {
	tt := []struct {
		expected []byte
		x        uint32
	}{
		{
			expected: []byte{1, 0, 0, 0},
			x:        1,
		}, {
			expected: []byte{0, 1, 0, 0},
			x:        256,
		}, {
			expected: []byte{0, 0, 1, 0},
			x:        65536,
		}, {
			expected: []byte{0, 0, 0, 1},
			x:        16777216,
		},
	}

	for _, tc := range tt {
		t.Run(strconv.Itoa(int(tc.x)), func(t *testing.T) {
			b := make([]byte, 4)
			rye.Serialize.Uint32(b, tc.x)
			assert.Equal(t, tc.expected, b)
			x := rye.Deserialize.Uint32(tc.expected)
			assert.Equal(t, tc.x, x)
		})
	}
}

func TestSerializeUint64(t *testing.T) {
	type testCase struct {
		expected []byte
		x        uint64
	}
	tt := make([]testCase, 8)
	x := uint64(1)
	for i := range tt {
		b := make([]byte, 8)
		b[i] = 1
		tt[i] = testCase{
			expected: b,
			x:        x,
		}
		x *= 256
	}
	for _, tc := range tt {
		t.Run(strconv.Itoa(int(tc.x)), func(t *testing.T) {
			b := make([]byte, 8)
			rye.Serialize.Uint64(b, tc.x)
			assert.Equal(t, tc.expected, b)
			x := rye.Deserialize.Uint64(tc.expected)
			assert.Equal(t, tc.x, x)
		})
	}
}

func TestSerializeAny(t *testing.T) {
	tt := []struct {
		expected []byte
		v        any
	}{
		{
			expected: []byte{0x63, 0xc5, 0x54, 0x0, 0x0, 0x0, 0x0, 0x0},
			v:        uint(5555555),
		}, {
			expected: []byte{0x50, 0x6},
			v:        uint16(1616),
		}, {
			expected: []byte{0xa0, 0x56, 0xa9, 0xc0},
			v:        uint32(3232323232),
		}, {
			expected: []byte{0x40, 0xad, 0x5f, 0xbc, 0x2b, 0xb4, 0xf8, 0x8},
			v:        uint64(646464646464646464),
		}, {
			expected: []byte{0x15, 0x81, 0xe9, 0x7d, 0xf4, 0x10, 0x22, 0x11},
			v:        int(1234567890123456789),
		}, {
			expected: []byte{0xff},
			v:        int8(-1),
		}, {
			expected: []byte{0xb0, 0xf9},
			v:        int16(-1616),
		}, {
			expected: []byte{0x60, 0xc9, 0x12, 0xfe},
			v:        int32(-32323232),
		}, {
			expected: []byte{0xc0, 0x52, 0xa0, 0x43, 0xd4, 0x4b, 0x7, 0xf7},
			v:        int64(-646464646464646464),
		}, {
			expected: []byte{0x56, 0xe, 0x49, 0x40},
			v:        float32(3.1415),
		}, {
			expected: []byte{0x6f, 0x12, 0x83, 0xc0, 0xca, 0x21, 0x9, 0x40},
			v:        float64(3.1415),
		}, {
			expected: []byte{10},
			v:        byte(10),
		}, {
			expected: []byte{0x74, 0x68, 0x69, 0x73, 0x20, 0x69, 0x73, 0x20, 0x61, 0x20, 0x74, 0x65, 0x73, 0x74},
			v:        "this is a test",
		}, {
			expected: []byte{1},
			v:        true,
		}, {
			expected: []byte{0},
			v:        false,
		}, {
			expected: nil,
			v:        []int{1, 2, 3},
		},
	}

	for _, tc := range tt {
		t.Run(reflect.TypeOf(tc.v).String(), func(t *testing.T) {
			assert.Equal(t, tc.expected, rye.Serialize.Any(tc.v, nil))
		})
	}
}

func TestFloat64OrderedFuzz(t *testing.T) {
	ln := 1000
	fs := make([]float64, ln)
	b := make([]byte, 8)
	for i := 0; i < ln; i++ {
		rand.Read(b)
		fs[i] = rye.Deserialize.Float64(b)
		if math.IsNaN(fs[i]) {
			i--
		}
	}
	sort.Float64Slice(fs).Sort()

	enc := make([][]byte, ln)
	for i, f := range fs {
		enc[i] = make([]byte, 8)
		rye.Serialize.Float64Ordered(enc[i], f)
	}
	sort.Slice(enc, func(i, j int) bool {
		return bytes.Compare(enc[i], enc[j]) == -1
	})

	got := make([]float64, len(fs))
	for i, e := range enc {
		got[i] = rye.Deserialize.Float64Ordered(e)
	}

	assert.Equal(t, fs, got)
}

func TestInt16OrderedFuzz(t *testing.T) {
	ln := 1000
	is := make(slice.Slice[int16], ln)
	b := make([]byte, 2)
	for i := range is {
		rand.Read(b)
		is[i] = rye.Deserialize.Int16(b)
	}
	is.Sort(slice.LT[int16]())

	enc := make([][]byte, ln)
	for i, i16 := range is {
		enc[i] = make([]byte, 2)
		rye.Serialize.Int16Ordered(enc[i], i16)
	}
	sort.Slice(enc, func(i, j int) bool {
		return bytes.Compare(enc[i], enc[j]) == -1
	})

	got := make(slice.Slice[int16], ln)
	for i, e := range enc {
		got[i] = rye.Deserialize.Int16Ordered(e)
	}

	assert.Equal(t, is, got)
}

func TestInt32OrderedFuzz(t *testing.T) {
	ln := 1000
	is := make(slice.Slice[int32], ln)
	b := make([]byte, 4)
	for i := range is {
		rand.Read(b)
		is[i] = rye.Deserialize.Int32(b)
	}
	is.Sort(slice.LT[int32]())

	enc := make([][]byte, ln)
	for i, i32 := range is {
		enc[i] = make([]byte, 4)
		rye.Serialize.Int32Ordered(enc[i], i32)
	}
	sort.Slice(enc, func(i, j int) bool {
		return bytes.Compare(enc[i], enc[j]) == -1
	})

	got := make(slice.Slice[int32], ln)
	for i, e := range enc {
		got[i] = rye.Deserialize.Int32Ordered(e)
	}

	assert.Equal(t, is, got)
}

func TestInt64OrderedFuzz(t *testing.T) {
	ln := 1000
	is := make(slice.Slice[int64], ln)
	b := make([]byte, 8)
	for i := range is {
		rand.Read(b)
		is[i] = rye.Deserialize.Int64(b)
	}
	is.Sort(slice.LT[int64]())

	enc := make([][]byte, ln)
	for i, i64 := range is {
		enc[i] = make([]byte, 8)
		rye.Serialize.Int64Ordered(enc[i], i64)
	}
	sort.Slice(enc, func(i, j int) bool {
		return bytes.Compare(enc[i], enc[j]) == -1
	})

	got := make(slice.Slice[int64], ln)
	for i, e := range enc {
		got[i] = rye.Deserialize.Int64Ordered(e)
	}

	assert.Equal(t, is, got)
}
