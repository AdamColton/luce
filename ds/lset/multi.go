package lset

import "github.com/adamcolton/luce/ds/slice"

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
