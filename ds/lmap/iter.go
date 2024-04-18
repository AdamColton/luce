package lmap

import (
	"github.com/adamcolton/luce/util/liter"
)

// KVGen creates a Key-Value pair from an input. Generally this will be from
// an iterator
type KVGen[T any, Key comparable, Val any] func(t T, idx int) (Key, Val, bool)

// FromIter creates a map from an iterator and a generator function.
func FromIter[I any, Key comparable, Val any](in liter.Iter[I], fn KVGen[I, Key, Val]) map[Key]Val {
	return itr(0, in, fn)
}

func itr[I any, K comparable, V any](size int, in liter.Iter[I], fn func(I, int) (K, V, bool)) map[K]V {
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
