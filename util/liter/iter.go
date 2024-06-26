// Package liter provides an iterator. To correctly implement an iterator,
// it should be initilized in a valid state so that this for loop would
// visit all the values:
//
//	for t,done := i.Cur(); !done; t,done = i.Next(){...}
package liter

import "sync"

// Iter interface allows for a standard set of tools for iterating over a
// collection.
type Iter[T any] interface {
	Next() (t T, done bool)
	Cur() (t T, done bool)
	Done() bool
	Idx() int
}

// Starter is an optional interface that Iter can implement to return to the
// start of the iteration.
type Starter[T any] interface {
	Start() (t T, done bool)
}

// Seek calls fn sequentially for each value Iter returns until Done is true.
// This does not reset the iterator.
func Seek[T any](i Iter[T], fn func(t T) bool) Iter[T] {
	t, done := i.Cur()
	return seek(i, t, done, fn)
}

// For calls fn sequentially for each value Iter. This does not reset the
// iterator.
func For[T any](i Iter[T], fn func(t T)) {
	t, done := i.Cur()
	fr(i, t, done, fn)
}

// For calls fn sequentially for each value Iter. This does not reset the
// iterator.
func ForIdx[T any](i Iter[T], fn func(t T, idx int)) int {
	t, done := i.Cur()
	return frIdx(i, t, done, fn)
}

// Concurrent calls fn in a Go routine for each value Iter returns until Done is
// true. The returned WaitGroup will reach zero when all Go routines return.
// This does not reset the iterator.
func Concurrent[T any](i Iter[T], fn func(t T, idx int)) *sync.WaitGroup {
	t, done := i.Cur()
	return concurrent(i, t, done, fn)
}

// Channel creates a chan with size buf and places each value from Iter on the
// channel until Done is true. This does not reset the iterator. The channel
// is filled from a Go routine so all the values need to be consumed or the
// routine will never close.
func Channel[T any](i Iter[T], buf int) <-chan T {
	t, done := i.Cur()
	return channel(i, t, done, buf)
}

// Pop returns the current value of iterator and if it is not done, calls Next.
func Pop[T any](i Iter[T]) T {
	t, done := i.Cur()
	if !done {
		i.Next()
	}
	return t
}
