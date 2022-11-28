package slice

// Buffer is used to provide a slice for re-use avoiding excessive allocation.
// Methods return a Slice that reuses the Buffer's memory when it has enough
// capacity, otherwise a new slice is allocated and the Buffer is left alone.
type Buffer[T any] []T

// Empty returns a zero length Slice with at least capacity c. If the buffer
// has capacity c, it will be used otherwise a new one is created.
func (buf Buffer[T]) Empty(c int) Slice[T] {
	if cap(buf) >= c {
		return Slice[T](buf[:0])
	}
	return make([]T, 0, c)
}

// Slice returns a Slice with length c. If the buffer has capacity c, it will be
// used and its values are left as they were, otherwise a new zeroed slice is
// created.
func (buf Buffer[T]) Slice(c int) Slice[T] {
	if cap(buf) >= c {
		return Slice[T](buf[:c])
	}
	return make([]T, c)
}

// Zeros returns a Slice with length c with all values set to the zero value. If
// the buffer has capacity c, it will be used otherwise a new one is created.
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

// ReduceCapacity sets the capacity to c, if that is lower than the current
// capacity. This can be useful when splitting a buffer to prevent use of the
// first part of the buffer from overflowing into the second part. The length is
// reduced to c if it is greater.
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

// Split a buffer, returning two from one. The first will have length 0 and
// capacity c. The second will have the remainder. If the buffer does not have
// capacity c, a new slice with capacity c is returned and the buffer is
// returned unchanged as the second value.
func (buf Buffer[T]) Split(c int) (Slice[T], Buffer[T]) {
	if cap(buf) < c {
		return make([]T, 0, c), buf
	}
	return buf[:0].ReduceCapacity(c), buf[c:c]
}
