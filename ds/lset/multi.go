package lset

import "github.com/adamcolton/luce/ds/slice"

type Multi[T comparable] []*Set[T]

// Sort Multi from smallest to largest. This order is assumend for optimization
// of other methods, but is not necessary.
func (m Multi[T]) Sort() {
	slice.Less[*Set[T]](func(i, j *Set[T]) bool {
		return i.Len() < j.Len()
	}).Sort(m)
}
