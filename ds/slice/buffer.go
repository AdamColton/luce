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
