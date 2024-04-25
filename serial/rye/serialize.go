package rye

import (
	"math"
	"unsafe"

	"github.com/adamcolton/luce/ds/slice"
)

// S provides a namespace to attach the Serialize methods.
type S struct{}

// Serialize holds the Serialize methods.
var Serialize S

// Uint16 encodes x and writes it to b.
func (S) Uint16(b []byte, x uint16) {
	b[0] = byte(x)
	b[1] = byte(x >> 8)
}

// Int32 encodes x and writes it to b.
func (s S) Int16(b []byte, x int16) {
	s.Uint16(b, uint16(x))
}

func (s S) postOrdered(b []byte) {
	Reverse(b)
	b[0] ^= 128
}

// Int16Ordered encodes x and writes it to b so that ordering it with
// bytes.Compare will produce the same ordering as when ordered by the int16
// values.
func (s S) Int16Ordered(b []byte, x int16) {
	s.Int16(b, x)
	s.postOrdered(b[:2])
}

// Uint32 encodes x and writes it to b.
func (S) Uint32(b []byte, x uint32) {
	b[0] = byte(x)
	b[1] = byte(x >> 8)
	b[2] = byte(x >> 16)
	b[3] = byte(x >> 24)
}

// Int32 encodes x and writes it to b.
func (s S) Int32(b []byte, x int32) {
	s.Uint32(b, uint32(x))
}

// Int32Ordered encodes x and writes it to b so that ordering it with
// bytes.Compare will produce the same ordering as when ordered by the int32
// values.
func (s S) Int32Ordered(b []byte, x int32) {
	s.Int32(b, x)
	s.postOrdered(b[:4])
}

// Uint64 encodes x and writes it to b.
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

// Int64 encodes x and writes it to b.
func (s S) Int64(b []byte, x int64) {
	s.Uint64(b, uint64(x))
}

// Int64Ordered encodes x and writes it to b so that ordering it with
// bytes.Compare will produce the same ordering as when ordered by the int64
// values.
func (s S) Int64Ordered(b []byte, x int64) {
	s.Int64(b, x)
	s.postOrdered(b[:8])
}

// Float64 encodes f and writes it to b.
func (s S) Float64(b []byte, f float64) {
	a := *(*[8]byte)(unsafe.Pointer(&f))
	b[0] = a[0]
	b[1] = a[1]
	b[2] = a[2]
	b[3] = a[3]
	b[4] = a[4]
	b[5] = a[5]
	b[6] = a[6]
	b[7] = a[7]
}

// Float64Ordered encodes f and writes it to b so that ordering it with
// bytes.Compare will produce the same ordering as when ordered by the float64
// values.
func (s S) Float64Ordered(b []byte, f float64) {
	s.Float64(b, f)
	Reverse(b)
	if b[0]&128 == 128 {
		Inverse(b)
	} else {
		b[0] |= 128
	}
}

// Any takes any type that rye can serialize (any int, uint, float, string or bool)
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
