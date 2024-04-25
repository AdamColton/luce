package rye

import "unsafe"

// D provides a namespace to attach the Deserialize methods.
type D struct{}

// Deserialize holds the Deserialize methods.
var Deserialize D

// Uint16 decodes b as a uint16
func (D) Uint16(b []byte) uint16 {
	return uint16(b[1])<<8 + uint16(b[0])
}

// Int16 decodes b as an int16
func (d D) Int16(b []byte) int16 {
	return int16(d.Uint16(b))
}

func (d D) preOrdered(b []byte) {
	b[0] ^= 128
	Reverse(b)
}

// Int16Ordered decodes b as an ordered int16
func (d D) Int16Ordered(b []byte) int16 {
	d.preOrdered(b[:2])
	return d.Int16(b)
}

// Uint32 decodes b as a uint32
func (D) Uint32(b []byte) uint32 {
	return uint32(b[3])<<24 +
		uint32(b[2])<<16 +
		uint32(b[1])<<8 +
		uint32(b[0])
}

// Int32 decodes b as an int32
func (d D) Int32(b []byte) int32 {
	return int32(d.Uint32(b))
}

// Int32Ordered decodes b as an ordered int32
func (d D) Int32Ordered(b []byte) int32 {
	d.preOrdered(b[:4])
	return d.Int32(b)
}

// Int64 decodes b as an int64
func (d D) Int64(b []byte) int64 {
	return int64(d.Uint64(b))
}

// Int64Ordered decodes b as an ordered int64
func (d D) Int64Ordered(b []byte) int64 {
	d.preOrdered(b[:8])
	return d.Int64(b)
}

// Uint64 decodes b as a uint64
func (D) Uint64(b []byte) uint64 {
	return uint64(b[7])<<56 +
		uint64(b[6])<<48 +
		uint64(b[5])<<40 +
		uint64(b[4])<<32 +
		uint64(b[3])<<24 +
		uint64(b[2])<<16 +
		uint64(b[1])<<8 +
		uint64(b[0])
}

// Float64 decodes b as a float64
func (D) Float64(b []byte) float64 {
	return *(*float64)(unsafe.Pointer(&b[0]))
}

// Float64Ordered decodes b as an ordered float64
func (d D) Float64Ordered(b []byte) float64 {
	if b[0]&128 == 128 {
		b[0] &= 127
	} else {
		Inverse(b)
	}
	Reverse(b[:8])
	return d.Float64(b)
}
