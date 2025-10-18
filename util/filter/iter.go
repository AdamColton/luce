package filter

import "github.com/adamcolton/luce/util/liter"

// Nexter fulfills liter.NextFunc
func (f Filter[T]) Nexter(i liter.Iter[T]) liter.NextFunc[T] {
	return func() (t T, done bool) {
		for t, done = i.Cur(); !done && !f(t); t, done = i.Next() {
		}
		i.Next()
		return
	}
}

// Iter created from the Filter.
func (f Filter[T]) Iter(i liter.Iter[T]) liter.Wrapper[T] {
	return f.Nexter(i).Indexer()
}
