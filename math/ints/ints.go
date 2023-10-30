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

// DivUp returns a/b rounding up. If a and b have different signs and b does not
// divide a evenly, the result is one too large.
func DivUp[T constraints.Integer](a, b T) T {
	out := a / b
	if a%b != 0 {
		out++
	}
	return out
}

// DivDown returns a/b using Go's integer division, which rounds towards zero.
// For a and b with the same sign that is rounding down. Defining the desired
// behavior can be more explicit.
func DivDown[T constraints.Integer](a, b T) T {
	return a / b
}

// == projects.Code.luce.ints ==
// [ ] DivRound(a, b T) T

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

// GCD finds the greatest common divisor of a and b. It is intended for
// non-negative values, the sign of the result is not reliable for negative ones.
func GCD[T constraints.Integer](a, b T) T {
	gcd, _, _ := GCDX(a, b)
	return gcd
}

// GCDX implements the extended GCD algorithm. It returns the GCD of a and b,
// and the coefficients x and y where a*x + b*y equals the GCD (Bezout's
// identity).
func GCDX[T constraints.Integer](a, b T) (gcd, x, y T) {
	if a == 0 {
		return b, 0, 1
	}
	gcd, u, v := GCDX(b%a, a)

	x = v - (b/a)*u
	y = u

	return gcd, x, y
}

// LCM finds the least common multiple of a and b. It panics if both are 0.
func LCM[T constraints.Integer](a, b T) T {
	return (a / GCD(a, b)) * b
}
