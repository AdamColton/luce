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

func (f Filter[T]) AndN(fs ...Filter[T]) Filter[T] {
	return func(val T) bool {
		if !f(val) {
			return false
		}
		for _, fn := range fs {
			if !fn(val) {
				return false
			}
		}
		return true
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

// New is syntactic sugar to make filter declarations a bit cleaner.
func New[T any](fn func(T) bool) Filter[T] {
	return fn
}

// FirstIdx finds the index of the first value that passes the filter
func (f Filter[T]) FirstIdx(vals ...T) int {
	for i, val := range vals {
		if f(val) {
			return i
		}
	}
	return -1
}

// First returns the first value that passes the Filter and the index. If none
// pass, then idx will be -1 and t will be the default value.
func (f Filter[T]) First(vals ...T) (t T, idx int) {
	idx = f.FirstIdx(vals...)
	if idx > -1 {
		t = vals[idx]
	}
	return
}

// SliceInPlace reorders the slice so all the elements passing the filter are at
// the start of the slice and all elements failing the filter are at the end.
// It returns two subslices, the first for passing, the second for failing.
// No guarentees are made about the order of the subslices.
func (f Filter[T]) SliceInPlace(vals []T) (passing, failing slice.Slice[T]) {
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
