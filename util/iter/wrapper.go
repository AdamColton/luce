package iter

import (
	"reflect"

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
