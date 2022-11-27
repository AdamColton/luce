package slice

import (
	"reflect"
	"sync"

	"github.com/adamcolton/luce/util/iter"
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
