package slice

import (
	"sync"

	"github.com/adamcolton/luce/util/iter"
)

type Slice[T any] []T

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

// Keys returns the keys of a map as a slice
func Keys[K comparable, V any](m map[K]V) Slice[K] {
	out := make([]K, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// Vals returns the values of a map as a slice.
func Vals[K comparable, V any](m map[K]V) Slice[V] {
	out := make([]V, 0, len(m))
	for _, v := range m {
		out = append(out, v)
	}
	return out
}

// Unique returns a slice with all the unique elements of the slice passed in.
func Unique[T comparable](s []T) Slice[T] {
	set := make(map[T]struct{})
	for _, t := range s {
		set[t] = struct{}{}
	}
	return Keys(set)
}

func (s Slice[T]) Iter() iter.Wrapper[T] {
	return NewIter(s)
}

func (s Slice[T]) IterFactory() (i iter.Iter[T], t T, done bool) {
	i = NewIter(s)
	t, done = i.Cur()
	return
}

// ForAll runs a Go routine for each element in s, passing it into fn. A
// WaitGroup is returned that will finish when all Go routines return.
func (s Slice[T]) ForAll(fn func(idx int, t T)) *sync.WaitGroup {
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
