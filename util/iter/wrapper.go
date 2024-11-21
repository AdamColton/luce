package iter

import (
	"reflect"
	"sync"

	"github.com/adamcolton/luce/util/upgrade"
)

// Wrapper provides useful methods that can be applied to any List.
type Wrapper[T any] struct {
	Iter[T]
}

// Wrap a Iter. Also checks that the underlying Iter is not itself a Wrapper.
func Wrap[T any](i Iter[T]) Wrapper[T] {
	if w, ok := i.(Wrapper[T]); ok {
		return w
	}
	return Wrapper[T]{i}
}

// Upgrade fulfills upgrade.Upgrader. Checks if the underlying Iter fulfills the
// given Type.
func (w Wrapper[T]) Upgrade(t reflect.Type) interface{} {
	return upgrade.Wrapped(w.Iter, t)
}

// Seek calls fn sequentially for each value Iter returns until Done is true.
// This does not reset the iterator.
func (w Wrapper[T]) Seek(fn func(t T) bool) Iter[T] {
	t, done := w.Cur()
	return seek(w.Iter, t, done, fn)
}

// For calls fn sequentially for each value Iter. This does not reset the
// iterator.
func (w Wrapper[T]) For(fn func(t T)) {
	t, done := w.Cur()
	fr(w.Iter, t, done, fn)
}

// For calls fn sequentially for each value Iter. This does not reset the
// iterator.
func (w Wrapper[T]) ForIdx(fn func(t T, idx int)) int {
	t, done := w.Cur()
	return frIdx(w.Iter, t, done, fn)
}

// Concurrent calls fn in a Go routine for each value Iter returns until Done is
// true. The returned WaitGroup will reach zero when all Go routines return.
// This does not reset the iterator.
func (w Wrapper[T]) Concurrent(fn func(t T, idx int)) *sync.WaitGroup {
	t, done := w.Cur()
	return concurrent(w.Iter, t, done, fn)
}

// Channel creates a chan with size buf and places each value from Iter on the
// channel until Done is true. This does not reset the iterator.
func (w Wrapper[T]) Channel(buf int) <-chan T {
	t, done := w.Cur()
	return channel(w.Iter, t, done, buf)
}