package iter

// Factory creates an iterator.
type Factory[T any] func() (iter Iter[T], t T, done bool)

// Seek creates a new Iter from the factory and calls fn sequentially for each
// value Iter returns until Done is true.
func (f Factory[T]) Seek(fn func(t T) bool) Iter[T] {
	i, t, done := f()
	return seek(i, t, done, fn)
}

// ForIdx calls fn sequentially for each value Iter. This does not reset the
// iterator.
func (f Factory[T]) For(fn func(t T)) {
	i, t, done := f()
	fr(i, t, done, fn)
}

// ForIdx calls fn sequentially for each value Iter. This does not reset the
// iterator.
func (f Factory[T]) ForIdx(fn func(t T, idx int)) int {
	i, t, done := f()
	return frIdx(i, t, done, fn)
}
