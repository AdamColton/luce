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
