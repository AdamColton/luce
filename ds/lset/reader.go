package lset

import "github.com/adamcolton/luce/ds/slice"

// Reader contains just the methods on Set that do not mutate the set
type Reader[T comparable] interface {
	Contains(elem T) bool
	Len() int
	Copy() *Set[T]
	Each(fn func(t T, done *bool))
	Slice(buf []T) slice.Slice[T]
}
