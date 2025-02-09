package lmap

import (
	"cmp"
	"fmt"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/util/liter"
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

// WrapNew returns a wrapper from the the underlying Mapper.New method.
func (w Wrapper[K, V]) WrapNew() Wrapper[K, V] {
	return Wrap(w.Mapper.New())
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

type Less[T any] = slice.Less[T]

// SortKeys is a convenience function that returns the sorted keys. This is
// equivalent to calling m.Keys(nil).Sort(slice.LT[K]()). It assumes slice.LT
// for sorting and a nil buffer. If either of those assumtions are not true, use
// Keys and Sort explicitly.
func SortKeys[K cmp.Ordered, V any](m map[K]V) slice.Slice[K] {
	return New(m).Keys(nil).Sort(cmp.Less[K])
}

// SortKeys creates a sorted slice of the keys.
func (w Wrapper[K, V]) SortKeys(less Less[K], buf []K) slice.Slice[K] {
	return w.Keys(buf).Sort(less)
}

// EachKey invokes the EachFunc for every
func (w Wrapper[K, V]) EachKey(i liter.Iter[K], fn EachFunc[K, V]) {
	for cur, done := i.Cur(); !done; cur, done = i.Next() {
		v, found := w.Get(cur)
		if found {
			fn(cur, v, &done)
			if done {
				break
			}
		}
	}
}

// SortedEachKey is shorthand for invoking SortedKeys and then EachKey. It
// creates a sorted slice as intermediary product.
func (w Wrapper[K, V]) SortedEachKey(less Less[K], buf []K, fn EachFunc[K, V]) slice.Slice[K] {
	keys := w.SortKeys(less, buf)
	w.EachKey(keys.Iter(), fn)
	return keys
}

// KeyLessKP is a helper that converts a Less function on the Key type to
// a Less funciton on the KeyPair type.
func KeyLessKP[V, K any](less Less[K]) Less[KeyPair[K, V]] {
	return func(i, j KeyPair[K, V]) bool {
		return less(i.Key, j.Key)
	}
}

func (w Wrapper[K, V]) KeyLessKP(less Less[K]) Less[KeyPair[K, V]] {
	return KeyLessKP[V](less)
}

func ValLessKP[K, V any](less Less[V]) Less[KeyPair[K, V]] {
	return func(i, j KeyPair[K, V]) bool {
		return less(i.Val, j.Val)
	}
}

func (w Wrapper[K, V]) ValLessKP(less Less[V]) Less[KeyPair[K, V]] {
	return ValLessKP[K](less)
}

// Slice converts the underlying map to a slice of keyPairs. If less is
// provided, it will be sorted.
func (w Wrapper[K, V]) Slice(less Less[KeyPair[K, V]], buf []KeyPair[K, V]) slice.Slice[KeyPair[K, V]] {
	out := slice.NewBuffer(buf).Cap(w.Len())
	w.Each(func(key K, val V, done *bool) {
		out = append(out, KeyPair[K, V]{
			Key: key,
			Val: val,
		})
	})
	if less != nil {
		out.Sort(less)
	}
	return out
}
