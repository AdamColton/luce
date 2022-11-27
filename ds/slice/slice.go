package slice

import (
	"reflect"
	"sort"

	"github.com/adamcolton/luce/util/liter"
)

// Slice is a wrapper that provides helper methods
type Slice[T any] []T

// New is syntactic sugar to infer the type
func New[T any](s []T) Slice[T] {
	return s
}

// Clone a slice.
func (s Slice[T]) Clone() Slice[T] {
	out := make([]T, len(s))
	copy(out, s)
	return out
}

// Swaps two values in the slice.
func (s Slice[T]) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// AppendNotZero will append any values from ts that are not the zero
// value for the type. Particularly useful for appending not nil values.
func (s Slice[T]) AppendNotZero(ts ...T) []T {
	for _, t := range ts {
		v := reflect.ValueOf(t)
		if v.Kind() != reflect.Invalid && !v.IsZero() {
			s = append(s, t)
		}
	}
	return s
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

// Remove values at given indicies by swapping them with values from the end
// and truncating the slice. Values less than zero or greater than the length
// of the list are ignored. Note that idxs is reordered so if that is a slice
// passed in and the order is important, pass in a copy.
func (s Slice[T]) Remove(idxs ...int) Slice[T] {
	sort.Sort(sort.Reverse(sort.IntSlice(idxs)))
	ln := len(s)
	prev := ln
	// Depending on variations in the implementation there are two things that
	// can make this behave in unintended ways. Duplicate values cause a double
	// swap. And it could be possible for a value near the end of the list to
	// removed, but then swapped with a value earlier in the list, reintroducing
	// it. Also, negative values are not allowed.
	//
	// To avoid both, idxs is sorted in descending order and prev tracks the
	// the last value. The "idx < prev" comparison guarentees both that there
	// are no duplicates and that idx is less than the length of the list.
	for _, idx := range idxs {
		if idx >= 0 && idx < prev {
			ln--
			s.Swap(idx, ln)
			prev = idx
		}
	}
	return s[:ln]
}
