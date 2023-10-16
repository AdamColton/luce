package liter

import (
	"github.com/adamcolton/luce/math/cmpr"
	"golang.org/x/exp/constraints"
)

// Reducer aggregates a value against every element in the iterator. It is passed
// the aggregate so far and returns the new aggregate.
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
// value. For instance if T is a struct, fn could return one of the fields. The
// aggregate passed in should be lower than any value fn returns, for example
// ints.MinI.
func Max[N constraints.Ordered, T any](fn func(T) N) Reducer[N, T] {
	return func(max N, element T, idx int) N {
		return cmpr.Max(max, fn(element))
	}
}

// Min value in the iter. The fn argument is used to convert to an ordered
// value. For instance if T is a struct, fn could return one of the fields. The
// aggregate passed in should be higher than any value fn returns, for example
// ints.MaxI.
func Min[N constraints.Ordered, T any](fn func(T) N) Reducer[N, T] {
	return func(max N, element T, idx int) N {
		return cmpr.Min(max, fn(element))
	}
}
