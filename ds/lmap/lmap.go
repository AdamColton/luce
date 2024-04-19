package lmap

import "github.com/adamcolton/luce/ds/slice"

type Map[K comparable, V any] map[K]V

func New[K comparable, V any](m map[K]V) Map[K, V] {
	return Map[K, V](m)
}

func (m Map[K, V]) Len() int {
	return len(m)
}

func (m Map[K, V]) Pop(key K) (V, bool) {
	v, found := m[key]
	if found {
		delete(m, key)
	}
	return v, found
}

func (m Map[K, V]) Vals(buf slice.Slice[V]) slice.Slice[V] {
	out := slice.NewBuffer(buf).Cap(len(m))
	for _, v := range m {
		out = append(out, v)
	}
	return out
}

func (m Map[K, V]) Keys(buf slice.Slice[K]) slice.Slice[K] {
	out := slice.NewBuffer(buf).Cap(len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
