package morph

import (
	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/util/liter"
	"github.com/adamcolton/luce/util/upgrade"
)

// Lener is fulfilled by types that know their length. It is the same as
// slice.Lener.
type Lener interface {
	Len() int
}

// Eacher applies the KeyValAll to every pair in e and returns an Eacher over the
// results. The keys of the result are the index of each result. If e fulfills
// Lener (possibly through a Wrapper) so does the result. Nothing is evaluated
// until Each is called.
func (kvt KeyValAll[K, V, Out]) Eacher(e Eacher[K, V]) Eacher[int, Out] {
	ln := 0
	if l, ok := upgrade.To[Lener](e); ok {
		ln = l.Len()
	}
	return liter.Eacher[Out]{
		L: ln,
		Func: func(inner EachFn[int, Out]) {
			idx := 0
			e.Each(func(k K, v V, done *bool) {
				out := kvt(k, v)
				inner(idx, out, done)
				idx++
			})
		},
	}
}

// KVtoKV is a func that converts a key and a value to a new key and value.
type KVtoKV[KIn, VIn, KOut, VOut any] = func(k KIn, v VIn) (KOut, VOut)

// NewKVToKV converts a KVtoKV to a KeyValAll that produces lmap.KeyVals, which
// lmap.FromEacher can use to create a map.
func NewKVToKV[KIn, VIn any, KOut comparable, VOut any](fn KVtoKV[KIn, VIn, KOut, VOut]) KeyValAll[KIn, VIn, lmap.KeyVal[KOut, VOut]] {
	return func(k KIn, v VIn) lmap.KeyVal[KOut, VOut] {
		return lmap.NewKV(fn(k, v))
	}
}

// OnVal creates a KeyValAll that keeps the key and applies the ValAll to the
// value.
func OnVal[K comparable, VIn, VOut any](fn ValAll[VIn, VOut]) KeyValAll[K, VIn, lmap.KeyVal[K, VOut]] {
	return func(k K, v VIn) lmap.KeyVal[K, VOut] {
		return lmap.NewKV(k, fn(v))
	}
}

// OnKey creates a KeyValAll that keeps the value and applies the ValAll to the
// key.
func OnKey[V any, KIn, KOut comparable](fn ValAll[KIn, KOut]) KeyValAll[KIn, V, lmap.KeyVal[KOut, V]] {
	return func(k KIn, v V) lmap.KeyVal[KOut, V] {
		return lmap.NewKV(fn(k), v)
	}
}

// Eacher applies the KeyVal to every pair in e and returns an Eacher over the
// results that are included. The keys of the result are the index of each
// included result. If e fulfills Lener (possibly through a Wrapper) so does the
// result, but its Len is the number of pairs in e, which is at least the number
// of included results. Nothing is evaluated until Each is called.
func (kvt KeyVal[K, V, Out]) Eacher(e Eacher[K, V]) Eacher[int, Out] {
	ln := 0
	if l, ok := upgrade.To[Lener](e); ok {
		ln = l.Len()
	}
	return liter.Eacher[Out]{
		L: ln,
		Func: func(inner EachFn[int, Out]) {
			idx := 0
			e.Each(func(k K, v V, done *bool) {
				out, include := kvt(k, v)
				if include {
					inner(idx, out, done)
					idx++
				}
			})
		},
	}
}
