package lset

import (
	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/slice"
)

type flag struct{}

// Set can contain values
type Set[T comparable] struct {
	m lmap.Map[T, flag]
}

// New creates a set containing the provided values.
func New[T comparable](elements ...T) *Set[T] {
	s := &Set[T]{
		m: make(lmap.Map[T, flag]),
	}
	s.Add(elements...)
	return s
}

// Contains return true if elem is in the set
func (s *Set[T]) Contains(elem T) bool {
	_, c := s.m[elem]
	return c
}

// Add given elements to the set
func (s *Set[T]) Add(elements ...T) {
	for _, t := range elements {
		s.m[t] = flag{}
	}
}

// Remove elem from the set
func (s *Set[T]) Remove(elem T) {
	delete(s.m, elem)
}

// Slice returns the values in the set as a slice
func (s *Set[T]) Slice() slice.Slice[T] {
	// TODO: take buf
	return s.m.Keys(nil)
}

// Len of the set
func (s *Set[T]) Len() int {
	return len(s.m)
}

// Copy the set
func (s *Set[T]) Copy() *Set[T] {
	out := &Set[T]{
		m: make(lmap.Map[T, flag], len(s.m)),
	}
	out.AddAll(s)
	return out
}

// AddAll elements of another set to this set
func (s *Set[T]) AddAll(set *Set[T]) {
	for k := range set.m {
		s.m[k] = flag{}
	}
}

// Each calls fn for each element in the set. This avoids the allocation of
// creating a slice when iterating over the values.
func (s *Set[T]) Each(fn func(T) (done bool)) {
	if s == nil {
		return
	}
	for t := range s.m {
		if done := fn(t); done {
			break
		}
	}
}
