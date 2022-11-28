package slice

type Buffer[T any] []T

// BufferEmpty returns a zero length buffer with at least capacity c. If the
// provided buffer has capacity, it will be used otherwise a new one is created.
func (buf Buffer[T]) Empty(c int) Slice[T] {
	if cap(buf) >= c {
		return Slice[T](buf[:0])
	}
	return make([]T, 0, c)
}

// BufferSlice returns a buffer with length c. If the provided buffer has
// capacity, it will be used otherwise a new one is created.
func (buf Buffer[T]) Slice(c int) Slice[T] {
	if cap(buf) >= c {
		return Slice[T](buf[:c])
	}
	return make([]T, c)
}
