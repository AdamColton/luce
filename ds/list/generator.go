package list

func NewGenerator[T any](ln int, fn func(int) T) Wrapper[T] {
	return Wrapper[T]{Generator[T]{
		Fn:     fn,
		Length: ln,
	}}
}

// NullaryGenerator uses a function that takes no arguments to create a list
// generator.
func NullaryGenerator[T any](ln int, fn func() T) Wrapper[T] {
	wrapped := func(idx int) T {
		return fn()
	}
	return NewGenerator(ln, wrapped)
}

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

// GeneratorFactory takes the generator function and returns a factory that only
// requires a length argument to create a Generator.
func GeneratorFactory[T any](fn func(int) T) func(ln int) Generator[T] {
	return func(ln int) Generator[T] {
		return Generator[T]{
			Fn:     fn,
			Length: ln,
		}
	}
}
