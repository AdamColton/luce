package lmap

import (
	"fmt"

	"github.com/adamcolton/luce/ds/slice"
	"golang.org/x/exp/constraints"
)

type Wrapper[K comparable, V any] struct {
	Mapper[K, V]
}

func Wrap[K comparable, V any](m Mapper[K, V]) Wrapper[K, V] {
	w, ok := m.(Wrapper[K, V])
	if ok {
		return w
	}
	return Wrapper[K, V]{m}
}

func (w Wrapper[K, V]) Wrapped() any {
	return w.Mapper
}

func (w Wrapper[K, V]) GetVal(key K) V {
	t, _ := w.Get(key)
	return t
}

func (w Wrapper[K, V]) Pop(key K) (V, bool) {
	v, found := w.Get(key)
	if found {
		w.Delete(key)
	}
	return v, found
}

func (w Wrapper[K, V]) MustPop(key K) V {
	v, ok := w.Pop(key)
	if !ok {
		panic(fmt.Errorf("failed to pop key: %v", key))
	}
	return v
}

func (w Wrapper[K, V]) Vals(buf slice.Slice[V]) slice.Slice[V] {
	out := slice.NewBuffer(buf).Cap(w.Len())
	w.Each(func(k K, v V, done *bool) {
		out = append(out, v)
	})
	return out
}

func (w Wrapper[K, V]) Keys(buf slice.Slice[K]) slice.Slice[K] {
	out := slice.NewBuffer(buf).Cap(w.Len())
	w.Each(func(k K, v V, done *bool) {
		out = append(out, k)
	})
	return out
}

// SortKeys is a convenience function that returns the keys sorted keys.
// This is equivalent to calling m.Keys(nil).Sort(slice.LT[K]()).
// It assumes slice.LT for sorting and a nil buffer. If either of those
// assumtions are not true, use Keys and Sort explicitly.
func SortKeys[K constraints.Ordered, V any](m map[K]V) slice.Slice[K] {
	return New(m).Keys(nil).Sort(slice.LT[K]())
}
