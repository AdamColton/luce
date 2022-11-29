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

func each[T any](i Iter[T], t T, iterDone bool, fn EachFn[T]) int {
	start := i.Idx()

	for fnDone := false; !iterDone && !fnDone; t, iterDone = i.Next() {
		fn(i.Idx(), t, &fnDone)
	}
	return i.Idx() - start
}

func concurrent[T any](i Iter[T], t T, iterDone bool, fn EachFn[T]) *sync.WaitGroup {
	if iterDone {
		return &sync.WaitGroup{}
	}
	fnDone := false
	mux := sync.Mutex{}
	return parallel.Run(func(coreIdx int) {
		innerDone := false
		for !iterDone && !fnDone {
			mux.Lock()
			if iterDone {
				mux.Unlock()
				return
			}
			lt := t
			idx := i.Idx()
			t, iterDone = i.Next()
			mux.Unlock()
			fn(idx, lt, &innerDone)
			if innerDone {
				fnDone = true
			}
		}
	})
}
