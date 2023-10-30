package ints

import "golang.org/x/exp/constraints"

// Provide Max value for all integer types and min value for all signed
// integers.
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

// DivUp returns a/b rounding up.
func DivUp[T constraints.Integer](a, b T) T {
	out := a / b
	if a%b != 0 {
		out++
	}
	return out
}

// DivDown returns a/b rounding down. This is the Go default, but defining the
// desired behavior can be more explicit.
func DivDown[T constraints.Integer](a, b T) T {
	return a / b
}

// Mod provides a version of modulus consistent with most other languages and
// calculators.
//
// In Go, mod (%) will return a negative if either a or b is negative. In most
// other languages and calculators the sign will always match b.
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

// GCD finds the greatest common denominator of a and b.
func GCD[T constraints.Integer](a, b T) T {
	gcd, _, _ := GCDX(a, b)
	return gcd
}

// GCDX implements the extended GCD algorithm. I do not understand what x and
// y represent.
func GCDX[T constraints.Integer](a, b T) (gcd, x, y T) {
	if a == 0 {
		return b, 0, 1
	}
	gcd, u, v := GCDX(b%a, a)

	x = v - (b/a)*u
	y = u

	return gcd, x, y
}

// LCM finds the least common multiple of a and b.
func LCM[T constraints.Integer](a, b T) T {
	return (a / GCD(a, b)) * b
}
