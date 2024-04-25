package rye

import (
	"math"

	"github.com/adamcolton/luce/ds/slice"
)

// TODO: why isn't this struct{}
type S byte

var Serialize S

// Uint16 writes a uint16 to the Serializer and increases the index.
func (S) Uint16(b []byte, x uint16) {
	b[0] = byte(x)
	b[1] = byte(x >> 8)
}

// Uint32 writes a uint32 to the Serializer and increases the index.
func (S) Uint32(b []byte, x uint32) {
	b[0] = byte(x)
	b[1] = byte(x >> 8)
	b[2] = byte(x >> 16)
	b[3] = byte(x >> 24)
}

// Uint64 writes a uint64 to the Serializer and increases the index.
func (S) Uint64(b []byte, x uint64) {
	b[0] = byte(x)
	b[1] = byte(x >> 8)
	b[2] = byte(x >> 16)
	b[3] = byte(x >> 24)
	b[4] = byte(x >> 32)
	b[5] = byte(x >> 40)
	b[6] = byte(x >> 48)
	b[7] = byte(x >> 56)
}

func (s S) Any(i any, buf []byte) []byte {
	switch t := i.(type) {
	case uint:
		buf = slice.NewBuffer(buf).Slice(8)
		s.Uint64(buf, uint64(t))
	case uint8:
		buf = slice.NewBuffer(buf).Slice(1)
		buf[0] = t
	case uint16:
		buf = slice.NewBuffer(buf).Slice(2)
		s.Uint16(buf, t)
	case uint32:
		buf = slice.NewBuffer(buf).Slice(4)
		s.Uint32(buf, t)
	case uint64:
		buf = slice.NewBuffer(buf).Slice(8)
		s.Uint64(buf, t)
	case int:
		buf = slice.NewBuffer(buf).Slice(8)
		s.Uint64(buf, uint64(t))
	case int8:
		buf = slice.NewBuffer(buf).Slice(1)
		buf[0] = byte(t)
	case int16:
		buf = slice.NewBuffer(buf).Slice(2)
		s.Uint16(buf, uint16(t))
	case int32:
		buf = slice.NewBuffer(buf).Slice(4)
		s.Uint32(buf, uint32(t))
	case int64:
		buf = slice.NewBuffer(buf).Slice(8)
		s.Uint64(buf, uint64(t))
	case float32:
		buf = slice.NewBuffer(buf).Slice(4)
		s.Uint32(buf, math.Float32bits(t))
	case float64:
		buf = slice.NewBuffer(buf).Slice(8)
		s.Uint64(buf, math.Float64bits(t))
	case string:
		buf = []byte(t)
	case bool:
		buf = slice.NewBuffer(buf).Slice(1)
		if t {
			buf[0] = 1
		} else {
			buf[0] = 0
		}
	default:
		buf = nil
	}
	return buf
}
