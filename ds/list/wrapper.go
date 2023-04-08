package list

import (
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/util/liter"
	"github.com/adamcolton/luce/util/upgrade"
)

// Wrapper provides a number of useful methods that can be applied to any List.
type Wrapper[T any] struct {
	List[T]
}

// Wrap a List. Also checks that the underlying list is not itself a Wrapper.
func Wrap[T any](l List[T]) Wrapper[T] {
	if w, ok := l.(Wrapper[T]); ok {
		return w
	}
	return Wrapper[T]{l}
}

// Wrapped fulfills upgrade.Wrapper.
func (w Wrapper[T]) Wrapped() any {
	return w.List
}

// Iter creates an liter.Iter backed by list L.
func (w Wrapper[T]) Iter() liter.Wrapper[T] {
	return NewIter(w.List)
}

// IterFactory creates an liter.Factory that generates a *list.Iter backed by
// list L.
func (w Wrapper[T]) IterFactory() liter.Factory[T] {
	return func() (it liter.Iter[T], t T, done bool) {
		it = &Iter[T]{
			List: w.List,
			I:    -1,
		}
		t, done = it.Next()
		return
	}
}

// Slice wraps a slice.Slice.
func Slice[T any](s []T) Wrapper[T] {
	return Wrap(slice.New(s))
}

// Reverse the list.
func (w Wrapper[T]) Reverse() Wrapper[T] {
	return Reverse[T](w).Wrap()
}

// Slice converts a List to slice. If the underlying List implements Slicer,
// that will be invoked.
func (w Wrapper[T]) Slice(buf []T) []T {
	if s, ok := upgrade.To[slice.Slicer[T]](w.List); ok {
		return s.Slice(buf)
	}
	return slice.FromIter(w.Iter(), buf)
}
