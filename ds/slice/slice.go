package slice

import (
	"reflect"
	"sort"

	"github.com/adamcolton/luce/math/cmpr"
	"github.com/adamcolton/luce/util/liter"
)

// Slice is a generic slice type that provides helper methods.
type Slice[T any] []T

// New converts s to a Slice, inferring the type.
func New[T any](s []T) Slice[T] {
	return s
}

// Make creates a Slice with the specified length and capacity. If cp is 0, ln
// is used for the capacity as well.
func Make[T any](ln, cp int) Slice[T] {
	if cp == 0 {
		cp = ln
	}
	return make(Slice[T], ln, cp)
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

// Remove values at given indices by swapping them with values from the end
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
	// last value. The "idx < prev" comparison guarantees both that there
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

// Buffer is syntactic sugar to convert a Slice to a Buffer providing a set of
// methods useful for buffering operations.
func (s Slice[T]) Buffer() Buffer[T] {
	return Buffer[T](s)
}

// Pop returns the last element of the slice and the slice resized to remove
// that element. If the size of the slice is zero, the zero value for the type
// is returned.
func (s Slice[T]) Pop() (T, Slice[T]) {
	ln := len(s)
	if ln == 0 {
		var t T
		return t, s
	}
	ln--
	return s[ln], s[:ln]
}

// Shift returns the first element of the slice and the slice resized to remove
// that element. If the size of the slice is zero, the zero value for the type
// is returned.
func (s Slice[T]) Shift() (T, Slice[T]) {
	ln := len(s)
	if ln == 0 {
		var t T
		return t, s
	}
	return s[0], s[1:ln]
}

// CheckCapacity ensures that Slice s has a capacity of at least check. If not,
// a new slice is created with capacity request and the slice is copied.
//
// If request is 0 (or any value less than check) then check is used. Request
// can be used to avoid additional allocations. For instance, doubling the
// capacity if there is insufficient space (as the underlying slice does in Go).
func (s Slice[T]) CheckCapacity(check, request int) Slice[T] {
	if cap(s) >= check {
		return s
	}
	request = cmpr.Max(check, request)
	out := make(Slice[T], len(s), request)
	copy(out, s)
	return out
}
