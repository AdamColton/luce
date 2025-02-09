package lmap

import (
	"cmp"
	"fmt"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/math/cmpr"
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

// Wrapped fulfills upgrade.Wrapper, returning the underlying Mapper.
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
// along with a bool indicating if the key was found.
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
// it. It will panic if the key is not found.
func (w Wrapper[K, V]) MustPop(key K) V {
	v, ok := w.Pop(key)
	if !ok {
		panic(fmt.Errorf("failed to pop key: %v", key))
	}
	return v
}

// All calls fn for every key/value pair. It is Each for a function that does
// not need done.
func (w Wrapper[K, V]) All(fn func(k K, v V)) {
	w.Each(All(fn))
}

// == projects.Code.luce.lmap ==
// [ ] Wrapper.(Vals,Keys) accept less

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

// WrapNew returns a Wrapper from the underlying Mapper.New method. The
// Mapper must not be nil.
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

// Copy returns a copy of the underlying map as a builtin map. The capacity can
// be set with cp, if it is less than the length of the map the length is used.
func (w Wrapper[K, V]) Copy(cp int) map[K]V {
	cp = cmpr.Max(w.Len(), cp)
	out := make(map[K]V, cp)
	w.Each(func(key K, val V, done *bool) {
		out[key] = val
	})
	return out
}

// SortKeys is a convenience function that returns the sorted keys. This is
// equivalent to calling m.Keys(nil).Sort(slice.LT[K]()). It assumes slice.LT
// for sorting and a nil buffer. If either of those assumptions are not true, use
// Keys and Sort explicitly.
func SortKeys[K cmp.Ordered, V any](m map[K]V) slice.Slice[K] {
	return New(m).Keys(nil).Sort(cmp.Less[K])
}

// Less is the same as slice.Less.
type Less[T any] = slice.Less[T]

// SortKeys creates a sorted slice of the keys.
func (w Wrapper[K, V]) SortKeys(less Less[K], buf []K) slice.Slice[K] {
	return w.Keys(buf).Sort(less)
}

// EachKey calls fn for every key the iterator provides, in the order it
// provides them, starting from the iterator's current value. Keys that are not
// in the map are skipped. Setting done in fn stops the iteration.
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

// SortedEachKey is shorthand for invoking SortKeys and then EachKey. It
// creates a sorted slice as intermediary product.
func (w Wrapper[K, V]) SortedEachKey(less Less[K], buf []K, fn EachFunc[K, V]) slice.Slice[K] {
	keys := w.SortKeys(less, buf)
	w.EachKey(keys.Iter(), fn)
	return keys
}

// KeyLessKP is a helper that converts a Less function on the Key type to
// a Less function on the KeyPair type.
func KeyLessKP[V, K any](less Less[K]) Less[KeyPair[K, V]] {
	return func(i, j KeyPair[K, V]) bool {
		return less(i.Key, j.Key)
	}
}

// KeyLessKP fulfills the same purpose as the KeyLessKP function, using the
// Wrapper's types.
func (w Wrapper[K, V]) KeyLessKP(less Less[K]) Less[KeyPair[K, V]] {
	return KeyLessKP[V](less)
}

// ValLessKP is a helper that converts a Less function on the Value type to a
// Less function on the KeyPair type.
func ValLessKP[K, V any](less Less[V]) Less[KeyPair[K, V]] {
	return func(i, j KeyPair[K, V]) bool {
		return less(i.Val, j.Val)
	}
}

// ValLessKP fulfills the same purpose as the ValLessKP function, using the
// Wrapper's types.
func (w Wrapper[K, V]) ValLessKP(less Less[V]) Less[KeyPair[K, V]] {
	return ValLessKP[K](less)
}

// Slice converts the underlying map to a slice of KeyPairs. The provided buffer
// will be used if it has sufficient capacity. If less is provided, the slice
// will be sorted; otherwise the order is not defined.
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
