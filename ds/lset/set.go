package lset

import "github.com/adamcolton/luce/ds/slice"

type flag struct{}

// Set can contain values
type Set[T comparable] struct {
	m map[T]flag
}

// New creates a set containing the provided values.
func New[T comparable](elements ...T) *Set[T] {
	s := &Set[T]{
		m: make(map[T]flag),
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
func (s *Set[T]) Slice() []T {
	return slice.Keys(s.m)
}

// Len of the set
func (s *Set[T]) Len() int {
	return len(s.m)
}

// Copy the set
func (s *Set[T]) Copy() *Set[T] {
	out := &Set[T]{
		m: make(map[T]flag, len(s.m)),
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
