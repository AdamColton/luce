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

func each[T any](i Iter[T], t T, iterDone bool, fn EachFn[T]) int {
	start := i.Idx()

	for fnDone := false; !iterDone && !fnDone; t, iterDone = i.Next() {
		fn(i.Idx(), t, &fnDone)
	}
	return i.Idx() - start
}
