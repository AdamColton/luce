package rye

import "github.com/adamcolton/luce/ds/slice"

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
	case uint64:
		buf = slice.NewBuffer(buf).Slice(8)
		s.Uint64(buf, t)
	case int:
		buf = slice.NewBuffer(buf).Slice(8)
		s.Uint64(buf, uint64(t))
	case int64:
		buf = slice.NewBuffer(buf).Slice(8)
		s.Uint64(buf, uint64(t))
	case string:
		buf = []byte(t)
	}
	return buf
}
