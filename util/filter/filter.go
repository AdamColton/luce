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

// SliceInPlace reorders the slice so all the elements passing the filter are at
// the start of the slice and all elements failing the filter are at the end.
// It returns two subslices, the first for passing, the second for failing.
// No guarentees are made about the order of the subslices.
func (f Filter[T]) SliceInPlace(vals []T) ([]T, []T) {
	ln := len(vals)
	if ln == 0 {
		return vals, nil
	}
	start := 0
	end := ln - 1
	for {
		for ; start < ln && f(vals[start]); start++ {
		}
		for ; end >= 0 && !f(vals[end]); end-- {
		}
		if start > end {
			break
		}
		vals[start], vals[end] = vals[end], vals[start]
	}
	return vals[:start], vals[start:]
}

// Chan runs a go routine listening on ch and any int that passes the Int is
// passed to the channel that is returned.
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

// Or builds a new Int that will return true if either underlying Int is true.
func (f Filter[T]) Or(f2 Filter[T]) Filter[T] {
	return func(val T) bool {
		return f(val) || f2(val)
	}
}

// And builds a new Int that will return true if both underlying Ints are true.
func (f Filter[T]) And(f2 Filter[T]) Filter[T] {
	return func(val T) bool {
		return f(val) && f2(val)
	}
}

// Not builds a new Int that will return true if the underlying Int is false.
func (f Filter[T]) Not() Filter[T] {
	return func(val T) bool {
		return !f(val)
	}
}

// Checker returns an error based on a single argument.
type Checker[T any] func(T) error

// Check converts a filter to a Checker and returns the provided err if the
// filter fails.
func (f Filter[T]) Check(errFn func(T) error) Checker[T] {
	return func(val T) error {
		if !f(val) {
			return errFn(val)
		}
		return nil
	}
}

// Panic runs the Checker and if it returns an error, panics with that error.
func (c Checker[T]) Panic(val T) {
	err := c(val)
	if err != nil {
		panic(err)
	}
}
