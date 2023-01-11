package slice

import (
	"reflect"
	"unsafe"

	"github.com/adamcolton/luce/util/upgrade"
)

// BufferEmpty returns a zero length buffer with at least capacity c. If the
// provided buffer has capacity, it will be used otherwise a new one is created.
func BufferEmpty[T any](c int, buf []T) []T {
	if cap(buf) >= c {
		return buf[:0]
	}
	return make([]T, 0, c)
}

// BufferLener returns a zero length buffer. If i fulfils Lener, the capacity of
// the buffer will be at least that returned by Len. If not, the buffer is used
// with the size set to zero.
func BufferLener[T any](i interface{}, buf []T) []T {
	ln := 0
	var ler Lener
	if upgrade.Upgrade(i, &ler) {
		ln = ler.Len()
	}
	return BufferEmpty(ln, buf)
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

// BufferSplit a buffer returns two buffers from one. The frist buffer will have
// capacity c. The second buffer will have the remainder. If the provided buffer
// does not have a capacity a new buffer is created.
func BufferSplit[T any](c int, buf []T) ([]T, []T) {
	if cap(buf) < c {
		return make([]T, 0, c), buf
	}
	return ReduceCapacity(c, buf[:0]), buf[c:c]
}
