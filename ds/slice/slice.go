package slice

import (
	"github.com/adamcolton/luce/math/cmpr"
	"github.com/adamcolton/luce/util/liter"
)

// Slice is a wrapper that provides helper methods
type Slice[T any] []T

// New is syntactic sugar to infer the type
func New[T any](s []T) Slice[T] {
	return s
}

// Clone a slice. The capacity can be set with cp. If cp is less than the length
// of s, that length will be used as the capacity.
func (s Slice[T]) Clone(cp int) Slice[T] {
	ln := len(s)
	cp = cmpr.Max(cp, ln)
	out := make([]T, ln, cp)
	copy(out, s)
	return out
}

// Swaps two values in the slice.
func (s Slice[T]) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Iter returns an iter.Wrapper for the slice.
func (s Slice[T]) Iter() liter.Wrapper[T] {
	return NewIter(s)
}

// IterFactory fulfills iter.Factory.
func (s Slice[T]) IterFactory() (i liter.Iter[T], t T, done bool) {
	i = NewIter(s)
	t, done = i.Cur()
	return
}
