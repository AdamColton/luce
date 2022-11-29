package liter

import (
	"sync"

	"github.com/adamcolton/luce/util/parallel"
)

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

func concurrent[T any](i Iter[T], t T, done bool, fn func(t T, idx int)) *sync.WaitGroup {
	if done {
		return &sync.WaitGroup{}
	}
	mux := sync.Mutex{}
	return parallel.Run(func(coreIdx int) {
		for !done {
			mux.Lock()
			if done {
				mux.Unlock()
				return
			}
			lt := t
			idx := i.Idx()
			t, done = i.Next()
			mux.Unlock()
			fn(lt, idx)
		}
	})
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
