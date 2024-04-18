package lmap

import (
	"github.com/adamcolton/luce/util/liter"
)

func Iter[K comparable, V any](in liter.Iter[K], fn func(K, int) (V, bool)) Map[K, V] {
	return itr(0, in, fn)
}

func itr[K comparable, V any](size int, in liter.Iter[K], fn func(K, int) (V, bool)) Map[K, V] {
	for i, done := in.Cur(); !done; i, done = in.Next() {
		o, include := fn(i, in.Idx())
		if include {
			in.Next()
			out := itr(size+1, in, fn)
			out[i] = o
			return out
		}
	}
	if size == 0 {
		return nil
	}
	return make(map[K]V, size)
}
