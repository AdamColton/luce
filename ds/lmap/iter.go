package lmap

import (
	"github.com/adamcolton/luce/util/liter"
)

// KeyPair holds a key and a value.
type KeyPair[K, V any] = struct {
	Key K
	Val V
}

// KeyVal is a KeyPair with a comparable key, so it can be used as a map entry.
type KeyVal[K comparable, V any] = KeyPair[K, V]

// NewKV creates a KeyVal, inferring the types.
func NewKV[K comparable, V any](k K, v V) KeyVal[K, V] {
	return KeyVal[K, V]{
		Key: k,
		Val: v,
	}
}

// FromIter creates a Wrapped Map from an iterator of KeyVals, starting from the
// iterator's current value. If the iterator is done, the Wrapper's Mapper is
// nil. If a key is repeated, the first value is kept. The values are collected
// recursively, so the stack grows with the length of the iterator.
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
