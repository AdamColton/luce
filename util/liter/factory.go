package liter

// Factory creates an iterator. It also returns the iterator's current value and
// done, as Cur would.
type Factory[T any] func() (iter Iter[T], t T, done bool)
