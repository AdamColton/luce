package filter

// Filter provides tools to filter ints and compose filters
type Filter[T any] func(T) bool

// Or builds a new Int that will return true if either underlying
// Int is true.
func (f Filter[T]) Or(f2 Filter[T]) Filter[T] {
	return func(val T) bool {
		return f(val) || f2(val)
	}
}

// And builds a new Int that will return true if both underlying
// Ints are true.
func (f Filter[T]) And(f2 Filter[T]) Filter[T] {
	return func(val T) bool {
		return f(val) && f2(val)
	}
}

// Not builds a new Int that will return true if the underlying
// Int is false.
func (f Filter[T]) Not() Filter[T] {
	return func(val T) bool {
		return !f(val)
	}
}
