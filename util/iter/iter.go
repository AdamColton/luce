package iter

import "sync"

// Iter interface allows for a standard set of tools for iterating over a
// collection.
type Iter[T any] interface {
	Next() (t T, done bool)
	Cur() (t T, done bool)
	Done() bool
	Idx() int
}

// Seek calls fn sequentially for each value Iter returns until Done is true.
// This does not reset the iterator.
func Seek[T any](i Iter[T], fn func(t T) bool) Iter[T] {
	t, done := i.Cur()
	return seek(i, t, done, fn)
}

// For calls fn sequentially for each value Iter. This does not reset the
// iterator.
func For[T any](i Iter[T], fn func(t T, idx int)) {
	t, done := i.Cur()
	idx := i.Idx()
	fr(i, t, done, idx, fn)
}

// Concurrent calls fn in a Go routine for each value Iter returns until Done is
// true. The returned WaitGroup will reach zero when all Go routines return.
// This does not reset the iterator.
func Concurrent[T any](i Iter[T], fn func(t T, idx int)) *sync.WaitGroup {
	t, done := i.Cur()
	return concurrent(i, t, done, fn)
}

// Channel creates a chan with size buf and places each value from Iter on the
// channel until Done is true. This does not reset the iterator.
func Channel[T any](i Iter[T], buf int) <-chan T {
	t, done := i.Cur()
	return channel(i, t, done, buf)
}
