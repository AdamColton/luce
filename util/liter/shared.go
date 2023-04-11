package liter

func seek[T any](i Iter[T], t T, done bool, fn func(t T) bool) Iter[T] {
	for ; !done; t, done = i.Next() {
		if fn(t) {
			return i
		}
	}
	return nil
}

func fr[T any](i Iter[T], t T, done bool, fn func(t T)) {
	for ; !done; t, done = i.Next() {
		fn(t)
	}
}

func frIdx[T any](i Iter[T], t T, done bool, fn func(t T, idx int)) int {
	start := i.Idx()
	for ; !done; t, done = i.Next() {
		fn(t, i.Idx())
	}
	return i.Idx() - start
}
