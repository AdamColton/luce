package liter

import "sync"

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

// For creates a new Iter from the factory and calls fn sequentially for each
// value Iter.
func (f Factory[T]) For(fn func(t T)) {
	i, t, done := f()
	fr(i, t, done, fn)
}

// Each calls fn sequentially for each value Iter.
func (f Factory[T]) Each(fn EachFn[T]) {
	i, t, done := f()
	each(i, t, done, fn)
}

// Concurrent creates a new Iter from the factory and calls fn in a Go routine
// for each value Iter returns until Done is true. The returned WaitGroup will
// reach zero when all Go routines return.
func (f Factory[T]) Concurrent(fn EachFn[T]) *sync.WaitGroup {
	i, t, done := f()
	return concurrent(i, t, done, fn)
}
