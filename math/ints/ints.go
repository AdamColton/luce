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

// LCMN finds the least common multiple of all integers in ns. It returns 0 if
// ns is empty.
func LCMN[T constraints.Integer](ns ...T) T {
	return Reduce(LCM, ns)
}

// ProdFn wraps a*b as a function
func ProdFn[T Number](a, b T) T {
	return a * b
}

// SumFn wraps a+b as a function.
func SumFn[T Number](a, b T) T {
	return a + b
}

// Prod returns the product of all the numbers in ns. It returns 1 if ns is
// empty.
func Prod[T Number](ns ...T) T {
	if len(ns) == 0 {
		return 1
	}
	return Reduce(ProdFn, ns)
}

// Range limits x to the range from start to end. If it is less than start,
// start is returned. If it is greater than end, end is returned.
func Range[T Number](start, x, end T) T {
	if x < start {
		return start
	}
	if x > end {
		return end
	}
	return x
}

// Int converts any integer type to an int. This and the other conversion
// functions truncate a value that does not fit in the result type, as a Go
// conversion does.
func Int[T constraints.Integer](i T) int { return int(i) }

// Int8 converts any integer type to an int8
func Int8[T constraints.Integer](i T) int8 { return int8(i) }

// Int16 converts any integer type to an int16
func Int16[T constraints.Integer](i T) int16 { return int16(i) }

// Int32 converts any integer type to an int32
func Int32[T constraints.Integer](i T) int32 { return int32(i) }

// Int64 converts any integer type to an int64
func Int64[T constraints.Integer](i T) int64 { return int64(i) }

// Uint converts any integer type to a uint
func Uint[T constraints.Integer](i T) uint { return uint(i) }

// Uint8 converts any integer type to a uint8
func Uint8[T constraints.Integer](i T) uint8 { return uint8(i) }

// Uint16 converts any integer type to a uint16
func Uint16[T constraints.Integer](i T) uint16 { return uint16(i) }

// Uint32 converts any integer type to a uint32
func Uint32[T constraints.Integer](i T) uint32 { return uint32(i) }

// Uint64 converts any integer type to a uint64
func Uint64[T constraints.Integer](i T) uint64 { return uint64(i) }

// Float32 converts any integer type to a float32
func Float32[T constraints.Integer](i T) float32 { return float32(i) }

// Float64 converts any integer type to a float64
func Float64[T constraints.Integer](i T) float64 { return float64(i) }

// Number is any integer or float.
type Number interface {
	constraints.Integer | constraints.Float
}

// Reducer is a function that combines two values into one. It is used to reduce
// a slice of values to a single value.
type Reducer[T Number] func(T, T) T

// Reduce combines the values in ts with the Reducer, from the first to the
// last. If ts is empty the zero value is returned.
func (fn Reducer[T]) Reduce(ts []T) (t T) {
	// Feels like this belongs in ds/slice or ds/list, but both of these
	// use math/ints and it creates an import cycle. This is also why it's
	// extended to constraints.Float
	if len(ts) == 0 {
		return
	}
	t = ts[0]
	for _, ti := range ts[1:] {
		t = fn(t, ti)
	}
	return
}

// Reduce calls fn on the first two elements of the slice, then for every
// element after that calls the function with the previous result as the first
// argument and the element as the second. If the slice is empty the zero value
// is returned.
func Reduce[T Number](fn Reducer[T], ts []T) (t T) {
	return fn.Reduce(ts)
}
