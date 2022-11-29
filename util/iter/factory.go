package iter

// Factory creates an iterator.
type Factory[T any] func() (iter Iter[T], t T, done bool)

// Do creates a new Iter from the factory and calls fn sequentially for each
// value Iter returns until Done is true.
func (f Factory[T]) Do(fn func(t T) bool) Iter[T] {
	i, t, done := f()
	return do(i, t, done, fn)
}
