package rye_test

import (
	"reflect"
	"strconv"
	"testing"

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
