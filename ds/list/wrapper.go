package list

import (
	"reflect"

	"github.com/adamcolton/luce/util/iter"
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

// Upgrade fulfills upgrade.Upgrader. Checks if the underlying List fulfills the
// given Type.
func (w Wrapper[T]) Upgrade(t reflect.Type) interface{} {
	return upgrade.Wrapped(w.List, t)
}

// Iter creates an iter.Iter backed by list L.
func (w Wrapper[T]) Iter() iter.Wrapper[T] {
	return iter.Wrapper[T]{&Iter[T]{List: w.List}}
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
