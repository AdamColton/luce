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
func (w Wrapper[K, V]) GetVal(key K) (v V) {
	if w.Mapper != nil {
		v, _ = w.Get(key)
	}
	return
}

// Each adds a nil check before calling Mapper.Each.
func (w Wrapper[K, V]) Each(fn EachFunc[K, V]) {
	if w.Mapper == nil {
		return
	}
	w.Mapper.Each(fn)
}

// Len adds a nil check before calling Mapper.Len.
func (w Wrapper[K, V]) Len() int {
	if w.Mapper == nil {
		return 0
	}
	return w.Mapper.Len()
}

// Pop removes a key from the map and returns the value associated with it
// along a with a bool indicating if the key was found.
func (w Wrapper[K, V]) Pop(key K) (v V, found bool) {
	if w.Mapper == nil {
		return
	}
	v, found = w.Get(key)
	if found {
		w.Delete(key)
	}
	return
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

// DeleteMany deletes multiple keys.
func (w Wrapper[K, V]) DeleteMany(keys []K) {
	if w.Mapper == nil {
		return
	}
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
