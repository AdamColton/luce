package lmap

import (
	"fmt"

	"github.com/adamcolton/luce/ds/slice"
	"golang.org/x/exp/constraints"
)

// Wrapper provides helpers around a Mapper.
type Wrapper[K comparable, V any] struct {
	Mapper[K, V]
}

// Wrap a Mapper. If it is already a Wrapper, that will be returned.
func Wrap[K comparable, V any](m Mapper[K, V]) Wrapper[K, V] {
	w, ok := m.(Wrapper[K, V])
	if ok {
		return w
	}
	return Wrapper[K, V]{m}
}

// Wrapped fulfills upgrade.Wrapper.
func (w Wrapper[K, V]) Wrapped() any {
	return w.Mapper
}

// GetVal returns the value for a key dropping the "found" boolean.
func (w Wrapper[K, V]) GetVal(key K) V {
	t, _ := w.Get(key)
	return t
}

// Pop removes a key from the map and returns the value associated with it
// along a with a bool indicating if the key was found.
func (w Wrapper[K, V]) Pop(key K) (V, bool) {
	v, found := w.Get(key)
	if found {
		w.Delete(key)
	}
	return v, found
}

// MustPop removes a key from the map and returns the value associated with
// it. It will panic nad the key is not found.
func (w Wrapper[K, V]) MustPop(key K) V {
	v, ok := w.Pop(key)
	if !ok {
		panic(fmt.Errorf("failed to pop key: %v", key))
	}
	return v
}

func (w Wrapper[K, V]) All(fn func(k K, v V)) {
	w.Each(All(fn))
}

// Vals returns the values of the map as a Slice. The provided buffer will be
// used if it has sufficient capacity.
func (w Wrapper[K, V]) Vals(buf slice.Slice[V]) slice.Slice[V] {
	out := slice.NewBuffer(buf).Cap(w.Len())
	w.Each(func(k K, v V, done *bool) {
		out = append(out, v)
	})
	return out
}

// Keys returns the keys of the map as a Slice. The provided buffer will be
// used if it has sufficient capacity.
func (w Wrapper[K, V]) Keys(buf slice.Slice[K]) slice.Slice[K] {
	out := slice.NewBuffer(buf).Cap(w.Len())
	w.Each(func(k K, v V, done *bool) {
		out = append(out, k)
	})
	return out
}

// WrapNew returns a wrapper from the the underlying Mapper.New method.
func (w Wrapper[K, V]) WrapNew() Wrapper[K, V] {
	return Wrap(w.Mapper.New())
}

// DeleteMany deletes multiple keys.
func (w Wrapper[K, V]) DeleteMany(keys []K) {
	for _, k := range keys {
		w.Mapper.Delete(k)
	}
}

// SortKeys is a convenience function that returns the sorted keys. This is
// equivalent to calling m.Keys(nil).Sort(slice.LT[K]()). It assumes slice.LT
// for sorting and a nil buffer. If either of those assumtions are not true, use
// Keys and Sort explicitly.
func SortKeys[K constraints.Ordered, V any](m map[K]V) slice.Slice[K] {
	return New(m).Keys(nil).Sort(slice.LT[K]())
}
