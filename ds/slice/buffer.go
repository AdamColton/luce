package slice

import (
	"github.com/adamcolton/luce/util/upgrade"
)

// Buffer is used to provide a slice for re-use avoiding excessive allocation.
type Buffer[T any] []T

// NewBuffer provides syntactic sugar by inferring the type.
func NewBuffer[T any](buf []T) Buffer[T] {
	return buf
}

// BufferEmpty returns a zero length buffer with at least capacity c. If the
// provided buffer has capacity, it will be used otherwise a new one is created.
func (buf Buffer[T]) Empty(c int) Slice[T] {
	if cap(buf) >= c {
		return Slice[T](buf[:0])
	}
	return make([]T, 0, c)
}

// Lener returns a zero length buffer. If i fulfils Lener, the capacity of the
// buffer will be at least that returned by Len. If not, the buffer is used with
// the size set to zero.
func (buf Buffer[T]) Lener(i any) Slice[T] {
	ln := 0
	ler, ok := upgrade.To[Lener](i)
	if ok {
		ln = ler.Len()
	}
	return buf.Empty(ln)
}

// BufferSlice returns a buffer with length c. If the provided buffer has
// capacity, it will be used otherwise a new one is created.
func (buf Buffer[T]) Slice(c int) Slice[T] {
	if cap(buf) >= c {
		return Slice[T](buf[:c])
	}
	return make([]T, c)
}

// Cap returns a buffer with at least capacity c and length 0. If the provided
// buffer has capacity, it will be used otherwise a new one is created.
func (buf Buffer[T]) Cap(c int) Slice[T] {
	if cap(buf) >= c {
		return Slice[T](buf[:0])
	}
	return make([]T, 0, c)
}

// BufferZeros returns a buffer with length c with all values set to 0. If the
// provided buffer has capacity, it will be used otherwise a new one is created.
func (buf Buffer[T]) Zeros(c int) Slice[T] {
	if cap(buf) >= c {
		var zero T
		buf = buf[:c]
		for i := range buf {
			buf[i] = zero
		}
		return Slice[T](buf)
	}
	return make([]T, c)
}

// ReduceCapacity sets the capacity to a lower value. This can be useful when
// splitting a buffer to prevent use of the first part of the buffer from
// overflowing into the second part.
func (buf Buffer[T]) ReduceCapacity(c int) Slice[T] {
	if c < cap(buf) {
		ln := len(buf)
		if c < ln {
			ln = c
		}
		return Slice[T](buf[:ln:c])
	}
	return Slice[T](buf)
}

// BufferSplit a buffer returns two buffers from one. The frist buffer will have
// capacity c. The second buffer will have the remainder. If the provided buffer
// does not have a capacity a new buffer is created.
func (buf Buffer[T]) Split(c int) (Slice[T], Buffer[T]) {
	if cap(buf) < c {
		return make([]T, 0, c), buf
	}
	return buf[:0].ReduceCapacity(c), buf[c:c]
}
