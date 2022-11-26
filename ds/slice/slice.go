package slice

import (
	"github.com/adamcolton/luce/math/cmpr"
	"github.com/adamcolton/luce/util/liter"
)

// Slice is a generic slice type that provides helper methods.
type Slice[T any] []T

// New converts s to a Slice, inferring the type.
func New[T any](s []T) Slice[T] {
	return s
}

// Clone returns a copy of the slice. The capacity can be set with cp. If cp is less than the length
// of s, that length will be used as the capacity. If cp is less than zero,
// then the length of s will be used.
func (s Slice[T]) Clone(cp int) Slice[T] {
	ln := len(s)
	if cp < 0 {
		cp = ln
	} else {
		cp = cmpr.Max(cp, ln)
	}
	out := make([]T, ln, cp)
	copy(out, s)
	return out
}

// Swap swaps two values in the slice.
func (s Slice[T]) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Iter returns a liter.Wrapper for the slice.
func (s Slice[T]) Iter() liter.Wrapper[T] {
	return NewIter(s)
}
