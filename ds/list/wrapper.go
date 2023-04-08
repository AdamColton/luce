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

// Wrap a List. If l is already a Wrapper, it is returned as is.
func Wrap[T any](l List[T]) Wrapper[T] {
	if w, ok := l.(Wrapper[T]); ok {
		return w
	}
	return Wrapper[T]{l}
}

// Wrapped fulfills upgrade.Wrapper, returning the underlying List.
func (w Wrapper[T]) Wrapped() any {
	return w.List
}

// Iter creates a liter.Wrapper that iterates over the List.
func (w Wrapper[T]) Iter() liter.Wrapper[T] {
	return NewIter(w.List)
}

// IterFactory creates a liter.Factory. Each call to the Factory creates a new
// *Iter over the List and returns it with its first value.
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

// Slice creates a Wrapper for the slice s.
func Slice[T any](s []T) Wrapper[T] {
	return Wrap(slice.New(s))
}

// Reverse returns a Wrapper for the List in reverse order. The underlying List
// is not changed.
func (w Wrapper[T]) Reverse() Wrapper[T] {
	return Reverse[T](w).Wrap()
}

// Slice converts a List to a slice. If the underlying List (possibly through a
// Wrapper) fulfills slice.Slicer, that will be invoked. Otherwise the values are
// copied into buf if it has enough capacity, or a new slice if not.
func (w Wrapper[T]) Slice(buf []T) []T {
	if s, ok := upgrade.To[slice.Slicer[T]](w.List); ok {
		return s.Slice(buf)
	}
	return slice.FromIter(w.Iter(), buf)
}
