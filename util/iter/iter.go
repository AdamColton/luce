package iter

import "sync"

// Iter interface allows for a standard set of tools for iterating over a
// collection.
type Iter[T any] interface {
	Start() (t T, done bool)
	Next() (t T, done bool)
	Cur() (t T, done bool)
	Done() bool
	Idx() int
}

// Do calls fn sequentially for each value Iter returns until Done is true. This
// does not reset the iterator.
func Do[T any](i Iter[T], fn func(t T) bool) Iter[T] {
	t, done := i.Cur()
	return do(i, t, done, fn)
}

// Concurrent calls fn in a Go routine for each value Iter returns until Done is
// true. The returned WaitGroup will reach zero when all Go routines return.
// This does not reset the iterator.
func Concurrent[T any](i Iter[T], fn func(t T, idx int)) *sync.WaitGroup {
	t, done := i.Cur()
	idx := i.Idx()
	return concurrent(i, t, done, idx, fn)
}
