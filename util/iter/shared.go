package iter

func seek[T any](i Iter[T], t T, done bool, fn func(t T) bool) Iter[T] {
	for ; !done; t, done = i.Next() {
		if fn(t) {
			return i
		}
	}
	return nil
}

func fr[T any](i Iter[T], t T, done bool, idx int, fn func(t T, idx int)) {
	for ; !done; t, done = i.Next() {
		fn(t, idx)
		idx++
	}
}
