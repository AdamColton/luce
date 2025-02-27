package lset

import (
	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/slice"
)

// Multi combines multiple sets and treats them as a single set
type Multi[T comparable] []*Set[T]

// NewMulti is a helper that infers type when creating a Multi
func NewMulti[T comparable](ts ...*Set[T]) Multi[T] {
	return ts
}

// Sort Multi from smallest to largest. This order is assumend for optimization
// of other methods, but is not necessary.
func (m Multi[T]) Sort() {
	slice.Less[*Set[T]](func(i, j *Set[T]) bool {
		return i.Len() < j.Len()
	}).Sort(m)
}

// Contains returns true if any Set in Multi contains t. The slice is checked
// in reverse order under the assumption that the slice is sorted from smallest
// to largest.
func (m Multi[T]) Contains(t T) bool {
	for idx := len(m) - 1; idx >= 0; idx-- {
		if m[idx].Contains(t) {
			return true
		}
	}
	return false
}

// AllContain returns true if every Set in Multi contains t. The slice is
// checked in order under the assumption that the slice is sorted from smallest
// to largest.
func (m Multi[T]) AllContain(t T) bool {
	for _, s := range m {
		if !s.Contains(t) {
			return false
		}
	}
	return true
}

// Intersection returns a set containing any value present in all sets.
func (m Multi[T]) Intersection() *Set[T] {
	if len(m) == 0 {
		return nil
	}
	if len(m) == 1 {
		return m[0].Copy()
	}
	out := &Set[T]{
		m: lmap.Empty[T, flag](m[0].Len()),
	}
	m1 := m[1:]
	m[0].m.Each(func(key T, val flag, done *bool) {
		if m1.AllContain(key) {
			out.m.Set(key, flag{})
		}
	})
	return out
}
