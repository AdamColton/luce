package bimap

import (
	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/slice"
)

// BiSlice pairs a slice and map. It can be particularly useful when populating
// a slice with unique values.
type BiSlice[T comparable] struct {
	Lookup lmap.Wrapper[T, int]
	Slice  slice.Slice[T]
}

// NewBiSlice creates a BiSlice. Either argument can be nil and a default will
// be created, but passing in a predefined values allow for memory efficiency
// or the Mapper can be a threadsafe variant.
func NewBiSlice[T comparable](m lmap.Mapper[T, int], buf []T) *BiSlice[T] {
	if m == nil {
		m = lmap.New(map[T]int{})
	}
	return &BiSlice[T]{
		Lookup: lmap.Wrapper[T, int]{m},
		Slice:  buf[:0],
	}
}

// Upser adds t to the slice and returns the index of where it is in the slice.
func (s *BiSlice[T]) Upsert(t T) int {
	idx, found := s.Lookup.Get(t)
	if !found {
		idx = len(s.Slice)
		s.Lookup.Set(t, idx)
		s.Slice = append(s.Slice, t)
	}
	return idx
}
