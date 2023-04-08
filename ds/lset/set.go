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
		m: lmap.New[T, flag](nil),
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
func (s *Set[T]) Add(elements ...T) {
	for _, t := range elements {
		s.m.Set(t, flag{})
	}
}

// Remove elem from the set
func (s *Set[T]) Remove(elem T) {
	s.m.Delete(elem)
}

// Slice returns the values in the set as a slice
func (s *Set[T]) Slice() slice.Slice[T] {
	// TODO: take buf
	return s.m.Keys(nil)
}

// Len of the set
func (s *Set[T]) Len() int {
	return s.m.Len()
}

// Copy the set
func (s *Set[T]) Copy() *Set[T] {
	out := &Set[T]{
		m: lmap.Empty[T, flag](s.m.Len()),
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
