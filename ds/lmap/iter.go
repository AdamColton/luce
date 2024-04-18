package lmap

import (
	"github.com/adamcolton/luce/util/liter"
)

type KeyPair[K, V any] = struct {
	Key K
	Val V
}

type KeyVal[K comparable, V any] = KeyPair[K, V]

func NewKV[K comparable, V any](k K, v V) KeyVal[K, V] {
	return KeyVal[K, V]{
		Key: k,
		Val: v,
	}
}

func FromIter[K comparable, V any](i liter.Iter[KeyVal[K, V]]) (out Wrapper[K, V]) {
	m := fromIter(i)
	if m != nil {
		out.Mapper = Map[K, V](m)
	}
	return
}

func fromIter[K comparable, V any](i liter.Iter[KeyVal[K, V]]) map[K]V {
	kv, done := i.Cur()
	if done {
		return nil
	}
	size := 1
	out := itr2(&size, i)
	out[kv.Key] = kv.Val
	return out
}

func itr2[K comparable, V any](size *int, i liter.Iter[KeyVal[K, V]]) map[K]V {
	kv, done := i.Next()
	if done {
		return make(map[K]V, *size)
	}
	*size++
	out := itr2(size, i)
	out[kv.Key] = kv.Val
	return out
}
