package liter

// Factory creates an iterator.
type Factory[T any] func() (iter Iter[T], t T, done bool)
