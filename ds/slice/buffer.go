package slice

import (
	"reflect"
	"unsafe"
)

// BufferEmpty returns a zero length buffer with at least capacity c. If the
// provided buffer has capacity, it will be used otherwise a new one is created.
func BufferEmpty[T any](c int, buf []T) []T {
	if cap(buf) >= c {
		return buf[:0]
	}
	return make([]T, 0, c)
}

// BufferSlice returns a buffer with length c. If the provided buffer has
// capacity, it will be used otherwise a new one is created.
func BufferSlice[T any](c int, buf []T) []T {
	if cap(buf) >= c {
		return buf[:c]
	}
	return make([]T, c)
}

// BufferZeros returns a buffer with length c with all values set to 0. If the
// provided buffer has capacity, it will be used otherwise a new one is created.
func BufferZeros[T any](c int, buf []T) []T {
	if cap(buf) >= c {
		var zero T
		buf = buf[:c]
		for i := range buf {
			buf[i] = zero
		}
		return buf
	}
	return make([]T, c)
}

// ReduceCapacity sets the capacity to a lower value. This can be useful when
// splitting a buffer to prevent use of the first part of the buffer from
// overflowing into the second part.
func ReduceCapacity[T any](c int, buf []T) []T {
	if c < cap(buf) {
		pv := reflect.ValueOf(&buf)
		sh := (*reflect.SliceHeader)(unsafe.Pointer(pv.Pointer()))
		sh.Cap = c
	}
	return buf
}
