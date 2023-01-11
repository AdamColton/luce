package list

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
func (w Wrapper[T]) Wrapped() any {
	return w.List
}
