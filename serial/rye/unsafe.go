package rye

import (
	"unsafe"
)

func float64ToUint64(f float64) uint64 {
	return *(*uint64)(unsafe.Pointer(&f))
}

func float32ToUint32(f float32) uint32 {
	return *(*uint32)(unsafe.Pointer(&f))
}

func uint32ToFloat32(u uint32) float32 {
	return *(*float32)(unsafe.Pointer(&u))
}

func uint64ToFloat64(u uint64) float64 {
	return *(*float64)(unsafe.Pointer(&u))
}
