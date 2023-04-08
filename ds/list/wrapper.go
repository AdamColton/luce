package list

import (
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/util/iter"
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

// Iter creates an iter.Iter backed by list L.
func (w Wrapper[T]) Iter() iter.Wrapper[T] {
	return NewIter(w.List)
}

// IterFactory creates an iter.Factory that generates a *list.Iter backed by
// list L.
func (w Wrapper[T]) IterFactory() iter.Factory[T] {
	return func() (it iter.Iter[T], t T, done bool) {
		it = &Iter[T]{
			List: w.List,
			I:    -1,
		}
		t, done = it.Next()
		return
	}
}

// Reverse the list.
func (w Wrapper[T]) Reverse() Wrapper[T] {
	return Reverse[T](w).Wrap()
}

// ToSlice converts a List to slice. If the underlying List implements Slicer,
// that will be invoked.
func (w Wrapper[T]) ToSlice(buf []T) []T {
	return slice.IterSlice[T](w.Iter(), buf)
}
