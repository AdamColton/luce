package filter

import (
	"github.com/adamcolton/luce/ds/channel"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
)

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

// Chan runs a go routine listening on ch and any int that passes the Int is
// passed to the channel that is returned.
func (f Filter[T]) Chan(pipe channel.Pipe[T]) channel.Pipe[T] {
	var out channel.Pipe[T]
	pipe, out.Snd, out.Rcv = channel.NewPipe(pipe.Rcv, pipe.Snd)
	go func() {
		for in := range pipe.Rcv {
			if f(in) {
				pipe.Snd <- in
			}
		}
		close(pipe.Snd)
	}()
	return out
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
	lerr.Panic(c(val))
}
