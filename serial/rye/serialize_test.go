package rye_test

import (
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
		})
	}
}
