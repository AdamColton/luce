package filter

import "github.com/adamcolton/luce/ds/slice"

// Filter represents boolean logic on a Type.
type Filter[T any] func(T) bool

// Or builds a new Filter that will return true if either underlying
// Filter is true.
func (f Filter[T]) Or(f2 Filter[T]) Filter[T] {
	return func(val T) bool {
		return f(val) || f2(val)
	}
}

// And builds a new Filter that will return true if both underlying
// Filters are true.
func (f Filter[T]) And(f2 Filter[T]) Filter[T] {
	return func(val T) bool {
		return f(val) && f2(val)
	}
}

// Not builds a new Filter that will return true if the underlying
// Filter is false.
func (f Filter[T]) Not() Filter[T] {
	return func(val T) bool {
		return !f(val)
	}
}

// SliceTransformFunc creates a slice.TransformFunc that uses the filter for
// the bool return argument.
func (f Filter[T]) SliceTransformFunc() slice.TransformFunc[T, T] {
	return func(t T, idx int) (T, bool) {
		return t, f(t)
	}
}

// Slice creates a new slice holding all values that return true when passed to
// Filter.
func (f Filter[T]) Slice(vals []T) slice.Slice[T] {
	return slice.TransformSlice(vals, nil, f.SliceTransformFunc())
}
