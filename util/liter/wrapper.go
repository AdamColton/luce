package liter

// Wrapper provides useful methods that can be applied to any Iter.
type Wrapper[T any] struct {
	Iter[T]
}

// Wrap wraps i in a Wrapper. If i is already a Wrapper it is returned as it is.
func Wrap[T any](i Iter[T]) Wrapper[T] {
	if w, ok := i.(Wrapper[T]); ok {
		return w
	}
	return Wrapper[T]{i}
}

// Wrapped fulfills upgrade.Wrapper.
func (w Wrapper[T]) Wrapped() any {
	return w.Iter
}
