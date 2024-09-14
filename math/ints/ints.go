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

// LCMN finds the least common multiple of all integers in ns
func LCMN[T constraints.Integer](ns ...T) T {
	return Compound(LCM, ns)
}

// ProdFn wraps a*b as a function
func ProdFn[T constraints.Integer](a, b T) T {
	return a * b
}

// SumFn wraps a+b as a function.
func SumFn[T constraints.Integer](a, b T) T {
	return a + b
}

// Prod returns the product of all integers in ns
func Prod[T constraints.Integer](ns ...T) T {
	if len(ns) == 0 {
		return 1
	}
	return Compound(ProdFn, ns)
}

func Compound[T constraints.Integer](fn func(a, b T) T, ts []T) (t T) {
	if len(ts) == 0 {
		return
	}
	t = ts[0]
	for _, ti := range ts[1:] {
		t = fn(t, ti)
	}
	return
}

// Idx provides a relative index to the given length. So a value of -1 will
// return the ln-1. The bool indicates if the index is in range.
func Idx(idx, ln int) (int, bool) {
	if idx < 0 {
		idx = ln + idx
	}
	return idx, idx >= 0 && idx < ln
}

func Range[T constraints.Integer](start, x, end T) T {
	if x < start {
		return start
	}
	if x > end {
		return end
	}
	return x
}

// Int converts any integer type to an int
func Int[T constraints.Integer](i T) int { return int(i) }

// Int8 converts any integer type to an int8
func Int8[T constraints.Integer](i T) int8 { return int8(i) }

// Int16 converts any integer type to an int16
func Int16[T constraints.Integer](i T) int16 { return int16(i) }

// Int32 converts any integer type to an int32
func Int32[T constraints.Integer](i T) int32 { return int32(i) }

// Int64 converts any integer type to an int64
func Int64[T constraints.Integer](i T) int64 { return int64(i) }

// Uint converts any integer type to an uint
func Uint[T constraints.Integer](i T) uint { return uint(i) }

// Uint8 converts any integer type to an uint8
func Uint8[T constraints.Integer](i T) uint8 { return uint8(i) }

// Uint16 converts any integer type to an uint16
func Uint16[T constraints.Integer](i T) uint16 { return uint16(i) }

// Uint32 converts any integer type to an uint32
func Uint32[T constraints.Integer](i T) uint32 { return uint32(i) }

// Uint64 converts any integer type to an uint64
func Uint64[T constraints.Integer](i T) uint64 { return uint64(i) }

// Float32 converts any integer type to an float32
func Float32[T constraints.Integer](i T) float32 { return float32(i) }

// Float64 converts any integer type to an float64
func Float64[T constraints.Integer](i T) float64 { return float64(i) }
