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
