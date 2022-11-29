package iter

import "sync"

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

func concurrent[T any](i Iter[T], t T, done bool, fn func(t T, idx int)) *sync.WaitGroup {
	var wg sync.WaitGroup
	wrap := func(t T, idx int) {
		fn(t, idx)
		wg.Add(-1)
	}
	for ; !done; t, done = i.Next() {
		idx := i.Idx()
		wg.Add(1)
		go wrap(t, idx)
	}
	return &wg
}

func channel[T any](i Iter[T], t T, done bool, buf int) <-chan T {
	ch := make(chan T, buf)
	go func() {
		for ; !done; t, done = i.Next() {
			ch <- t
		}
		close(ch)
	}()
	return ch
}
