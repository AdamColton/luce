package iter

// Iter interface allows for a standard set of tools for iterating over a
// collection.
type Iter[T any] interface {
	Next() (t T, done bool)
	Cur() (t T, done bool)
	Done() bool
	Idx() int
}

// Seek calls fn sequentially for each value Iter returns until Done is true.
// This does not reset the iterator.
func Seek[T any](i Iter[T], fn func(t T) bool) Iter[T] {
	t, done := i.Cur()
	return seek(i, t, done, fn)
}
