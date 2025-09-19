package liter

// Nexter only requires a Next function and can be used by the Indexer to
// create an Iter.
type Nexter[T any] interface {
	Next() (t T, done bool)
}

// NextFunc fulfills the Nexter interface
type NextFunc[T any] func() (t T, done bool)

// Next invokes the underlying function and fulfills Nexter
func (fn NextFunc[T]) Next() (t T, done bool) {
	return fn()
}

// Indexer creates an indexer. It invokes Next once for the initial values.
func (fn NextFunc[T]) Indexer() Wrapper[T] {
	return IndexerFactory(fn.Factory)
}

// Factory invokes NextFunc.Next once to get the initial values. It fullfils
// NexterFactory.
func (fn NextFunc[T]) Factory() (Nexter[T], T, bool) {
	cur, done := fn()
	return fn, cur, done
}

// NewNextFunc is a helper to infer parameterized types.
func NewNextFunc[T any](fn NextFunc[T]) NextFunc[T] {
	return fn
}

// NexterFactory is a helper whose output can be passed directly into
// NewIndexer.
type NexterFactory[T any] func() (Nexter[T], T, bool)

// Indexer tracks the current value and index of the underlying Nexter.
type Indexer[T any] struct {
	n    Nexter[T]
	idx  int
	cur  T
	done bool
}

// NewIndexer creates a Wrapped Indexer from a Nexter. It fulfills Iter[T].
func NewIndexer[T any](n Nexter[T], cur T, done bool) Wrapper[T] {
	return Wrap(&Indexer[T]{
		n:    n,
		cur:  cur,
		done: done,
	})
}

// IndexerFactory creates a Wrapped Indexer from a NexterFactory. It fulfills
// Iter[T].
func IndexerFactory[T any](factory NexterFactory[T]) Wrapper[T] {
	return NewIndexer(factory())
}

// Next invokes the underlying Nexter storing the values for cur and done
// incrementing the index. It fulfills Iter[T].
func (i *Indexer[T]) Next() (t T, done bool) {
	i.cur, i.done = i.n.Next()
	i.idx++
	return i.Cur()
}

// Cur returns the values for cur and done incrementing the index. It fulfills
// Iter[T].
func (i *Indexer[T]) Cur() (t T, done bool) {
	return i.cur, i.done
}

// Done returns true if the underlying Nexter is done. It fulfills Iter[T].
func (i *Indexer[T]) Done() bool {
	return i.done
}

// Idx returns the index. It fulfills Iter[T].
func (i *Indexer[T]) Idx() int {
	return i.idx
}
