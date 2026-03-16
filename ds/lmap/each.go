package lmap

import "github.com/adamcolton/luce/util/upgrade"

// Each fulfills Eacher and Lener using a function.
type Each[K comparable, V any] struct {
	// Func is called by Each with the EachFunc.
	Func func(EachFunc[K, V])
	// L is returned by Len.
	L int
}

// Each calls Func.
func (e Each[K, V]) Each(fn EachFunc[K, V]) {
	e.Func(fn)
}

// Len returns L.
func (e Each[K, V]) Len() int {
	return e.L
}

// IdxEacher is fulfilled by any type that can call a function on each of a
// series of key/value pairs, indexed by position.
type IdxEacher[K comparable, V any] interface {
	Each(EachFunc[int, KeyVal[K, V]])
}

// Lener is fulfilled by types that know their length.
type Lener interface {
	Len() int
}

// FromEach creates a Wrapped Map from the key/value pairs the function passes
// to the EachFunc. If a key is repeated, the last value is kept.
func FromEach[K comparable, V any](fn func(EachFunc[int, KeyVal[K, V]])) Wrapper[K, V] {
	m := make(map[K]V)
	fn(func(idx int, kv KeyVal[K, V], done *bool) {
		m[kv.Key] = kv.Val
	})
	return Wrap(Map[K, V](m))
}

// FromEacher creates a Wrapped Map from an IdxEacher. If the IdxEacher fulfills
// Lener (possibly through a Wrapper), the length is used to size the map. If a
// key is repeated, the last value is kept.
func FromEacher[K comparable, V any](e IdxEacher[K, V]) Wrapper[K, V] {
	ln := 0
	if l, ok := upgrade.To[Lener](e); ok {
		ln = l.Len()
	}
	m := make(map[K]V, ln)
	e.Each(func(idx int, kv KeyVal[K, V], done *bool) {
		m[kv.Key] = kv.Val
	})
	return Wrap(Map[K, V](m))
}

// AppendEacher sets every pair from the IdxEacher on the Wrapper and returns
// the Wrapper. The Mapper must not be nil.
func (w Wrapper[K, V]) AppendEacher(e IdxEacher[K, V]) Wrapper[K, V] {
	return w.AppendEach(e.Each)
}

// AppendEach sets every pair passed to the EachFunc on the Wrapper and returns
// the Wrapper. The Mapper must not be nil.
func (w Wrapper[K, V]) AppendEach(fn func(EachFunc[int, KeyVal[K, V]])) Wrapper[K, V] {
	fn(func(idx int, kv KeyVal[K, V], done *bool) {
		w.Set(kv.Key, kv.Val)
	})
	return w
}
