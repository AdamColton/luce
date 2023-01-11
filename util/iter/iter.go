package iter

// Iter interface allows for a standard set of tools for iterating over a
// collection.
type Iter[T any] interface {
	Next() (t T, done bool)
	Cur() (t T, done bool)
	Done() bool
	Idx() int
}

// Starter is an optional interface that Iter can implement to return to the
// start of the iteration.
type Starter[T any] interface {
	Start() (t T, done bool)
}

// Seek calls fn sequentially for each value Iter returns until Done is true.
// This does not reset the iterator.
func Seek[T any](i Iter[T], fn func(t T) bool) Iter[T] {
	t, done := i.Cur()
	return seek(i, t, done, fn)
}

// For calls fn sequentially for each value Iter. This does not reset the
// iterator.
func For[T any](i Iter[T], fn func(t T)) {
	t, done := i.Cur()
	fr(i, t, done, fn)
}
