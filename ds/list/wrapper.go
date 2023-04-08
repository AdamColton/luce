package list

import (
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/util/liter"
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
