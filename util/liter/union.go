package liter

// Union is a slice of Iter where each call moves all the Iter in unison.
type Union[T any] []Iter[T]

// NewUnion create a Untion.
func NewUnion[T any](is ...Iter[T]) Wrapper[[]T] {
	return Wrapper[[]T]{Union[T](is)}
}

// Next returns the next value of each Iter in the slice. Done will be true if
// any iter is done
func (u Union[T]) Next() (t []T, done bool) {
	t = make([]T, 0, len(u))
	for _, i := range u {
		it, d := i.Next()
		done = done || d
		t = append(t, it)
	}
	return
}

// Cur returns the current value of each Iter in the slice. Done will be true if
// any iter is done
func (u Union[T]) Cur() (t []T, done bool) {
	t = make([]T, 0, len(u))
	for _, i := range u {
		it, d := i.Cur()
		done = done || d
		t = append(t, it)
	}
	return
}

// Done returns true if any Iter in the Union is done.
func (u Union[T]) Done() bool {
	if len(u) == 0 {
		return true
	}
	for _, i := range u {
		if i.Done() {
			return true
		}
	}
	return false
}

// Idx returns the Idx of the first Iter in the slice. The expectation is that
// each Iter in the Union has the same index.
func (u Union[T]) Idx() int {
	if len(u) == 0 {
		return 0
	}
	return u[0].Idx()
}
