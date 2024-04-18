package lmap

import (
	"github.com/adamcolton/luce/util/liter"
)

func FromIter[I any, K comparable, V any](in liter.Iter[I], fn func(I, int) (K, V, bool)) Wrapper[K, V] {
	return New(itr(0, in, fn))
}

func itr[I any, K comparable, V any](size int, in liter.Iter[I], fn func(I, int) (K, V, bool)) Map[K, V] {
	for i, done := in.Cur(); !done; i, done = in.Next() {
		k, v, include := fn(i, in.Idx())
		if include {
			in.Next()
			out := itr(size+1, in, fn)
			out[k] = v
			return out
		}
	}
	if size == 0 {
		return nil
	}
	return make(map[K]V, size)
}
