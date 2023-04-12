package liter

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
