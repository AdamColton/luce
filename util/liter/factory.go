package liter

import "sync"

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

// Concurrent creates a new Iter from the factory and calls fn in a Go routine
// for each value Iter returns until Done is true. The returned WaitGroup will
// reach zero when all Go routines return.
func (f Factory[T]) Concurrent(fn func(t T, idx int)) *sync.WaitGroup {
	i, t, done := f()
	return concurrent(i, t, done, fn)
}

// Channel creates a new Iter from the factory and creates a chan with size buf
// and places each value from Iter on the channel until Done is true. The
// channel is filled from a Go routine so all the values need to be consumed or
// the routine will never close.
func (f Factory[T]) Channel(buf int) <-chan T {
	i, t, done := f()
	return channel(i, t, done, buf)
}

// Wrap invokes the Factory and wraps the returned iterator.
func (f Factory[T]) Wrap() (iter Wrapper[T], t T, done bool) {
	iter.Iter, t, done = f()
	return
}
