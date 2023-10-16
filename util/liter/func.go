package liter

import (
	"github.com/adamcolton/luce/math/cmpr"
	"golang.org/x/exp/constraints"
)

// Reducer aggregates a value against every elemnt in the iterator.
type Reducer[A, T any] func(aggregate A, element T, idx int) A

func (r Reducer[A, T]) reduce(t T, done bool, idx int, aggregate A, i Iter[T]) A {
	for ; !done; t, done = i.Next() {
		aggregate = r(aggregate, t, idx)
		idx++
	}
	return aggregate
}

// Iter runs the Reducer against an Iterator.
func (r Reducer[A, T]) Iter(aggregate A, i Iter[T]) A {
	t, done := i.Cur()
	return r.reduce(t, done, i.Idx(), aggregate, i)
}

// Factory runs the Reducer against an Iterator generated from the given
// Factory.
func (r Reducer[A, T]) Factory(aggregate A, f Factory[T]) A {
	i, t, done := f()
	return r.reduce(t, done, i.Idx(), aggregate, i)
}

// Appender creates a reducer that appends to a slice.
func Appender[T any]() Reducer[[]T, T] {
	return func(aggregate []T, element T, idx int) []T {
		return append(aggregate, element)
	}
}

// Max value in the iter. The fn argument is used to convert to an ordered
// value. For instance if T is a struct, fn could return one of the fields.
func Max[N constraints.Ordered, T any](fn func(T) N) Reducer[N, T] {
	return func(max N, element T, idx int) N {
		return cmpr.Max(max, fn(element))
	}
}

// Min value in the iter. The fn argument is used to convert to an ordered
// value. For instance if T is a struct, fn could return one of the fields.
func Min[N constraints.Ordered, T any](fn func(T) N) Reducer[N, T] {
	return func(max N, element T, idx int) N {
		return cmpr.Min(max, fn(element))
	}
}
