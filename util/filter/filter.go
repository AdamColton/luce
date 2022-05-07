package filter

// Filter provides tools to filter ints and compose filters
type Filter[T any] func(T) bool

// Returns all values that return true when passed to Int.
func (f Filter[T]) Slice(vals []T) []T {
	var out []T
	for _, val := range vals {
		if f(val) {
			out = append(out, val)
		}
	}
	return out
}

// Chan runs a go routine listening on ch and any int that passes the
// Int is passed to the channel that is returned.
func (f Filter[T]) Chan(ch <-chan T, buf int) <-chan T {
	out := make(chan T, buf)
	go func() {
		for in := range ch {
			if f(in) {
				out <- in
			}
		}
		close(out)
	}()
	return out
}

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
