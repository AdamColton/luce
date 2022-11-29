package slice

import (
	"sort"
	"sync"
)

// Clone a slice.
func Clone[T any](s []T) []T {
	out := make([]T, len(s))
	copy(out, s)
	return out
}

// Swaps two values in the slice.
func Swap[T any](s []T, i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Keys returns the keys of a map as a slice
func Keys[K comparable, V any](m map[K]V) []K {
	out := make([]K, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// Vals returns the values of a map as a slice.
func Vals[K comparable, V any](m map[K]V) []V {
	out := make([]V, 0, len(m))
	for _, v := range m {
		out = append(out, v)
	}
	return out
}

// AppendNotDefault will append any values from ts that are not the default
// value for the type. Particularly useful for appending not nil values.
func AppendNotDefault[T comparable](s []T, ts ...T) []T {
	var d T
	for _, t := range ts {
		if d != t {
			s = append(s, t)
		}
	}
	return s
}

// Unique returns a slice with all the unique elements of the slice passed in.
func Unique[T comparable](s []T) []T {
	set := make(map[T]struct{})
	for _, t := range s {
		set[t] = struct{}{}
	}
	return Keys(set)
}

// ForAll runs a Go routine for each element in s, passing it into fn. A
// WaitGroup is returned that will finish when all Go routines return.
func ForAll[T any](s []T, fn func(idx int, t T)) *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(len(s))
	wrap := func(idx int, t T) {
		fn(idx, t)
		wg.Add(-1)
	}
	for i, t := range s {
		go wrap(i, t)
	}
	return &wg
}

// Remove values at given indicies by swapping them with values from the end
// and truncating the slice. Values less than zero or greater than the length
// of the list are ignored. Note that idxs is reordered so if that is a slice
// passed in and the order is important, pass in a copy.
func Remove[T any](s []T, idxs ...int) []T {
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
			Swap(s, idx, ln)
			prev = idx
		}
	}
	return s[:ln]
}

// Pop returns the last element of the slice and the slice resized to remove
// that element. If the size of the slice is zero, the zero value for the type
// is returned.
func Pop[T any](s []T) (T, []T) {
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
func Shift[T any](s []T) (T, []T) {
	ln := len(s)
	if ln == 0 {
		var t T
		return t, s
	}
	return s[0], s[1:ln]
}
