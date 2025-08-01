package lset

import "github.com/adamcolton/luce/ds/slice"

// Reader contains the methods on Set that do not mutate the set, other than All
// and SortedEach.
type Reader[T comparable] interface {
	Contains(elem T) bool
	Len() int
	Copy() *Set[T]
	Each(fn IterFunc[T])
	Slice(buf []T) slice.Slice[T]
}
