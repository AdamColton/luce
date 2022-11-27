package slice

import (
	"reflect"

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

// AppendNotZero appends any values from ts that are not the zero value for the
// type, and returns the result. Particularly useful for appending not nil
// values. A nil interface value counts as zero.
func (s Slice[T]) AppendNotZero(ts ...T) []T {
	for _, t := range ts {
		v := reflect.ValueOf(t)
		if v.Kind() != reflect.Invalid && !v.IsZero() {
			s = append(s, t)
		}
	}
	return s
}

// Iter returns a liter.Wrapper for the slice.
func (s Slice[T]) Iter() liter.Wrapper[T] {
	return NewIter(s)
}

// IterFactory fulfills liter.Factory. It returns a new iterator over the slice
// and its first value.
func (s Slice[T]) IterFactory() (i liter.Iter[T], t T, done bool) {
	i = NewIter(s)
	t, done = i.Cur()
	return
}
