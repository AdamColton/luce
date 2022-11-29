package liter

// Iter interface allows for a standard set of tools for iterating over a
// collection. To correctly implement an Iter, it should be initialized in a
// valid state so that this for loop would visit all the values:
//
//	for t, done := i.Cur(); !done; t, done = i.Next() {...}
type Iter[T any] interface {
	// Next moves to the next value and returns it. done is true once there are
	// no more values.
	Next() (t T, done bool)
	// Cur returns the current value without moving. done is true if there is no
	// current value because the iteration is finished or empty.
	Cur() (t T, done bool)
	// Done returns true if there are no more values.
	Done() bool
	// Idx returns the index of the current value.
	Idx() int
}

// Starter is an optional interface that Iter can implement to return to the
// start of the iteration.
type Starter[T any] interface {
	Start() (t T, done bool)
}

// Seek calls fn sequentially for each value Iter returns, starting with the
// current one, until fn returns true. It returns the Iter positioned on that
// value, or nil if the iterator is done first. This does not reset the
// iterator.
func Seek[T any](i Iter[T], fn func(t T) bool) Iter[T] {
	t, done := i.Cur()
	return seek(i, t, done, fn)
}
