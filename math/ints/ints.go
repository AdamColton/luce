package ints

import "golang.org/x/exp/constraints"

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

func DivUp[T constraints.Integer](a, b T) T {
	out := a / b
	if a%b != 0 {
		out++
	}
	return out
}

func DivDown[T constraints.Integer](a, b T) T {
	return a / b
}

func Mod[T constraints.Integer](a, b T) T {
	if a < 0 {
		return (b - (-a % b)) % b
	}
	return a % b
}

func GCD[T constraints.Integer](a, b T) T {
	gcd, _, _ := GCDX(a, b)
	return gcd
}

func GCDX[T constraints.Integer](a, b T) (gcd, q, r T) {
	if a == 0 {
		return b, 0, 1
	}
	gcd, u, v := GCDX(b%a, a)

	q = v - (b/a)*u
	r = u

	return gcd, q, r
}

func LCM[T constraints.Integer](a, b T) T {
	return (a / GCD(a, b)) * b
}

func LCMN[T constraints.Integer](ns ...T) T {
	return Compound(LCM, ns)
}

func ProdFn[T constraints.Integer](a, b T) T {
	return a * b
}

func SumFn[T constraints.Integer](a, b T) T {
	return a + b
}

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

func Int[T constraints.Integer](i T) int     { return int(i) }
func Int8[T constraints.Integer](i T) int8   { return int8(i) }
func Int16[T constraints.Integer](i T) int16 { return int16(i) }
func Int32[T constraints.Integer](i T) int32 { return int32(i) }
func Int64[T constraints.Integer](i T) int64 { return int64(i) }

func Uint[T constraints.Integer](i T) uint     { return uint(i) }
func Uint8[T constraints.Integer](i T) uint8   { return uint8(i) }
func Uint16[T constraints.Integer](i T) uint16 { return uint16(i) }
func Uint32[T constraints.Integer](i T) uint32 { return uint32(i) }
func Uint64[T constraints.Integer](i T) uint64 { return uint64(i) }

func Float32[T constraints.Integer](i T) float32 { return float32(i) }
func Float64[T constraints.Integer](i T) float64 { return float64(i) }
