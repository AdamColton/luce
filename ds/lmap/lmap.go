package lmap

import (
	"github.com/adamcolton/luce/ds/slice"
	"golang.org/x/exp/constraints"
)

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

// SortKeys is a convenience function that returns the keys sorted keys.
// This is equivalent to calling m.Keys(nil).Sort(slice.LT[K]()).
// It assumes slice.LT for sorting and a nil buffer. If either of those
// assumtions are not true, use Keys and Sort explicitly.
func SortKeys[K constraints.Ordered, V any](m map[K]V) slice.Slice[K] {
	return Map[K, V](m).Keys(nil).Sort(slice.LT[K]())
}
