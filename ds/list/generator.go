package list

// Generator fulfills List using a function to generate values by index.
type Generator[T any] struct {
	Fn     func(int) T
	Length int
}

// AtIdx fulfills List returning the value at the specified index.
func (g Generator[T]) AtIdx(idx int) T {
	return g.Fn(idx)
}

// Len fulfills List returning the length of the list.
func (g Generator[T]) Len() int {
	return g.Length
}

// Wrap the Generator to add Wrapper methods.
func (g Generator[T]) Wrap() Wrapper[T] {
	return Wrapper[T]{g}
}

func GeneratorFactory[T any](fn func(int) T) func(ln int) Generator[T] {
	return func(ln int) Generator[T] {
		return Generator[T]{
			Fn:     fn,
			Length: ln,
		}
	}
}
