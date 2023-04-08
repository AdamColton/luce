package lset

import "github.com/adamcolton/luce/ds/slice"

type flag struct{}

type Set[T comparable] struct {
	m map[T]flag
}

func New[T comparable](ts ...T) *Set[T] {
	s := &Set[T]{
		m: make(map[T]flag),
	}
	s.Add(ts...)
	return s
}

func (s *Set[T]) Contains(t T) bool {
	_, c := s.m[t]
	return c
}

func (s *Set[T]) Add(ts ...T) {
	for _, t := range ts {
		s.m[t] = flag{}
	}
}

func (s *Set[T]) Remove(t T) {
	delete(s.m, t)
}

func (s *Set[T]) Slice() []T {
	return slice.Keys(s.m)
}

func (s *Set[T]) Len() int {
	return len(s.m)
}

func (s *Set[T]) Copy() *Set[T] {
	out := &Set[T]{
		m: make(map[T]flag, len(s.m)),
	}
	out.AddAll(s)
	return out
}

func (s *Set[T]) AddAll(set *Set[T]) {
	for k := range set.m {
		s.m[k] = flag{}
	}
}
