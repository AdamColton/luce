package slice

import (
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
