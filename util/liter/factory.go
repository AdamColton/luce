package liter

// Factory creates an iterator. It also returns the iterator's current value and
// done, as Cur would.
type Factory[T any] func() (iter Iter[T], t T, done bool)

// Seek creates a new Iter from the factory and calls fn sequentially for each
// value Iter returns until fn returns true. It returns the Iter positioned on
// that value, or nil if the iterator is done first.
func (f Factory[T]) Seek(fn func(t T) bool) Iter[T] {
	i, t, done := f()
	return seek(i, t, done, fn)
}
