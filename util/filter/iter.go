package filter

import "github.com/adamcolton/luce/util/iter"

// Iter wraps an Iter and applies a Filter to it. Only values that pass the
// filter are returned.
type Iter[T any] struct {
	In iter.Iter[T]
	Filter[T]
	idx int
}

// Iter created from the Filter.
func (f Filter[T]) Iter(i iter.Iter[T]) iter.Wrapper[T] {
	for t, done := i.Cur(); !done && !f(t); t, done = i.Next() {
	}
	return iter.Wrap(&Iter[T]{
		In:     i,
		Filter: f,
	})
}

// Next fulfills iter.Iter, moves to the next value in the underlying Iter that
// passes the Filter or the default value of T if the iterator is done. Returns
// a bool indicating if iteration is done.
func (i *Iter[T]) Next() (t T, done bool) {
	for t, done = i.In.Next(); !done && !i.Filter(t); t, done = i.In.Next() {
	}
	i.idx++
	return
}

// Cur fulfills iter.Iter and returns the current value of the iterator and
// a bool indicating if iteration is done.
func (i *Iter[T]) Cur() (t T, done bool) {
	return i.In.Cur()
}

// Done returns a bool indicating if iteration is done.
func (i *Iter[T]) Done() bool {
	return i.In.Done()
}

// Idx returns the current index. This index is associated with the filtered
// iterator, not the underlying iterator.
func (i *Iter[T]) Idx() int {
	return i.idx
}
