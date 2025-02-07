package lset

import (
	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/slice"
)

type flag struct{}

// Set can contain values
type Set[T comparable] struct {
	m lmap.Wrapper[T, flag]
}

// New creates a set containing the provided values.
func New[T comparable](elements ...T) *Set[T] {
	s := &Set[T]{
		m: lmap.Empty[T, flag](len(elements)),
	}
	s.Add(elements...)
	return s
}

// Safe creates a threadsafe set
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

// Add given elements to the set
func (s *Set[T]) Add(elements ...T) *Set[T] {
	for _, t := range elements {
		s.m.Set(t, flag{})
	}
	return s
}

// Remove elem from the set
func (s *Set[T]) Remove(elem T) *Set[T] {
	s.m.Delete(elem)
	return s
}

// Slice returns the values in the set as a slice. If the Set is nil, a nil
// slice is returned.
func (s *Set[T]) Slice(buf []T) slice.Slice[T] {
	if s == nil {
		return nil
	}
	return s.m.Keys(buf)
}

// Len of the set
func (s *Set[T]) Len() int {
	return s.m.Len()
}

// Copy the set
func (s *Set[T]) Copy() *Set[T] {
	out := &Set[T]{
		m: s.m.WrapNew(),
	}
	out.AddAll(s)
	return out
}

// AddAll elements of another set to this set
func (s *Set[T]) AddAll(set *Set[T]) {
	set.m.Each(func(key T, val flag, done *bool) {
		s.m.Set(key, flag{})
	})
}

// Each calls fn for each element in the set. This avoids the allocation of
// creating a slice when iterating over the values.
func (s *Set[T]) Each(fn func(t T, done *bool)) {
	if s == nil {
		return
	}
	s.m.Each(func(key T, val flag, done *bool) {
		fn(key, done)
	})
}

// All calls the function for every element in the set.
func (s *Set[T]) All(fn func(t T)) {
	if s == nil {
		return
	}
	s.m.Each(func(key T, val flag, done *bool) {
		fn(key)
	})
}
