package ints

import "golang.org/x/exp/constraints"

// The Max constants give the largest value of each integer type and the Min
// constants give the smallest value of each signed integer type. The smallest
// value of an unsigned type is 0.
const (
	MaxU   uint   = ^uint(0)
	MaxU8  uint8  = ^uint8(0)
	MaxU16 uint16 = ^uint16(0)
	MaxU32 uint32 = ^uint32(0)
	MaxU64 uint64 = ^uint64(0)

	MaxI   int   = int(MaxU >> 1)
	MaxI8  int8  = int8(MaxU8 >> 1)
	MaxI16 int16 = int16(MaxU16 >> 1)
	MaxI32 int32 = int32(MaxU32 >> 1)
	MaxI64 int64 = int64(MaxU64 >> 1)

	MinI   int   = ^MaxI
	MinI8  int8  = ^MaxI8
	MinI16 int16 = ^MaxI16
	MinI32 int32 = ^MaxI32
	MinI64 int64 = ^MaxI64
)

// Dir is the direction Div rounds in when b does not divide a evenly.
type Dir byte

// The directions Div can round in.
const (
	// Down rounds toward negative infinity.
	Down Dir = iota
	// Up rounds toward positive infinity.
	Up
	// In rounds toward zero, the same as Go's integer division.
	In
	// Out rounds away from zero.
	Out
)

// Div returns a/b rounded in the direction d, whatever the signs of a and b.
// A d other than Down, Up, In or Out rounds like In. It panics if b is 0.
func Div[T constraints.Integer](a, b T, d Dir) T {
	q := a / b
	if a%b == 0 {
		return q
	}
	// Go has rounded q toward zero, which is up for a negative quotient and down
	// for a positive one.
	negative := (a < 0) != (b < 0)
	switch {
	case negative && (d == Down || d == Out):
		q--
	case !negative && (d == Up || d == Out):
		q++
	}
	return q
}

// DivDown returns a/b rounded toward negative infinity. It panics if b is 0.
func DivDown[T constraints.Integer](a, b T) T {
	return Div(a, b, Down)
}

// DivUp returns a/b rounded toward positive infinity. It panics if b is 0.
func DivUp[T constraints.Integer](a, b T) T {
	return Div(a, b, Up)
}

// DivIn returns a/b rounded toward zero, the same as Go's integer division. It
// panics if b is 0.
func DivIn[T constraints.Integer](a, b T) T {
	return Div(a, b, In)
}

// DivOut returns a/b rounded away from zero. It panics if b is 0.
func DivOut[T constraints.Integer](a, b T) T {
	return Div(a, b, Out)
}

// Mod provides a version of modulus consistent with most other languages and
// calculators.
//
// In Go, mod (%) will return a negative if either a or b is negative. In most
// other languages and calculators the sign will always match b. Mod panics if b
// is 0.
func Mod[T constraints.Integer](a, b T) T {
	if a < 0 {
		m := (b - (-a % b)) % b
		return m
	}
	m := a % b
	if m > 0 && b < 0 {
		m += b
	}
	return m
}
