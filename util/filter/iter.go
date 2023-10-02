package filter

import "github.com/adamcolton/luce/util/liter"

// Nexter returns a liter.NextFunc that returns the next value of i that passes
// the Filter, and done once i is finished.
func (f Filter[T]) Nexter(i liter.Iter[T]) liter.NextFunc[T] {
	return func() (t T, done bool) {
		for t, done = i.Cur(); !done && !f(t); t, done = i.Next() {
		}
		i.Next()
		return
	}
}

// Iter returns an iterator over the values of i that pass the Filter.
func (f Filter[T]) Iter(i liter.Iter[T]) liter.Wrapper[T] {
	return f.Nexter(i).Indexer()
}
