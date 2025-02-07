package lset

import (
	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/slice"
)

type flag struct{}

// Set is a set of comparable values. Create one with New or Safe. A nil *Set can
// be used with Len (0), Slice (nil), Each, All and SortedEach (nothing is
// called); any other method panics on a nil Set.
type Set[T comparable] struct {
	m lmap.Wrapper[T, flag]
}

// New creates a set containing the provided values. Repeated values are only
// in the set once.
func New[T comparable](elements ...T) *Set[T] {
	s := &Set[T]{
		m: lmap.Empty[T, flag](len(elements)),
	}
	s.Add(elements...)
	return s
}

// Safe creates a threadsafe set containing the provided values. Each, All and
// SortedEach hold a read lock while the func runs, so the func must not add to
// or remove from the set.
func Safe[T comparable](elements ...T) *Set[T] {
	s := &Set[T]{
		m: lmap.EmptySafe[T, flag](len(elements)),
	}
	s.Add(elements...)
	return s
}

// Contains return true if elem is in the set
func (s *Set[T]) Contains(elem T) bool {
	_, c := s.m.Get(elem)
	return c
}

// Add given elements to the set. The Set is returned to allow chaining.
func (s *Set[T]) Add(elements ...T) *Set[T] {
	for _, t := range elements {
		s.m.Set(t, flag{})
	}
	return s
}

// Remove elem from the set. It is not an error if elem is not in the set. The
// Set is returned to allow chaining.
func (s *Set[T]) Remove(elem T) *Set[T] {
	s.m.Delete(elem)
	return s
}

// Slice returns the values in the set as a slice, in no particular order. The
// buffer is used if it has sufficient capacity. If the Set is nil, a nil slice
// is returned.
func (s *Set[T]) Slice(buf []T) slice.Slice[T] {
	if s == nil {
		return nil
	}
	return s.m.Keys(buf)
}

// Len returns the number of elements in the set. If the Set is nil, it is 0.
func (s *Set[T]) Len() int {
	if s == nil {
		return 0
	}
	return s.m.Len()
}

// Copy the set into a new Set of the same kind, so a copy of a Safe set is safe.
// It panics if the Set is nil.
func (s *Set[T]) Copy() *Set[T] {
	out := &Set[T]{
		m: s.m.WrapNew(),
	}
	out.AddAll(s)
	return out
}

// AddAll elements of another set to this set. It panics if set is nil.
func (s *Set[T]) AddAll(set *Set[T]) {
	set.m.Each(func(key T, val flag, done *bool) {
		s.m.Set(key, flag{})
	})
}

// IterFunc is the func type for Each. Setting done to true stops the iteration.
type IterFunc[T any] func(t T, done *bool)

// Each calls fn for each element in the set, in no particular order. This avoids
// the allocation of creating a slice when iterating over the values. If the Set
// is nil, fn is not called.
func (s *Set[T]) Each(fn IterFunc[T]) {
	if s == nil {
		return
	}
	s.m.Each(func(key T, val flag, done *bool) {
		fn(key, done)
	})
}

// All calls the function for every element in the set, in no particular order.
// If the Set is nil, fn is not called.
func (s *Set[T]) All(fn func(t T)) {
	if s == nil {
		return
	}
	s.m.Each(func(key T, val flag, done *bool) {
		fn(key)
	})
}

// SortedEach first sorts the values with less, using buf if it has sufficient
// capacity, and calls fn with the values in that order. Setting done stops the
// iteration. If the Set is nil, fn is not called.
func (s *Set[T]) SortedEach(less slice.Less[T], buf []T, fn IterFunc[T]) {
	if s == nil {
		return
	}
	s.m.SortedEachKey(less, buf, func(key T, val flag, done *bool) {
		fn(key, done)
	})
}
