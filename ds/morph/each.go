package morph

import (
	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/util/liter"
	"github.com/adamcolton/luce/util/upgrade"
)

type Lener interface {
	Len() int
}

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

type KVtoKV[KIn, VIn, KOut, VOut any] = func(k KIn, v VIn) (KOut, VOut)

func NewKVToKV[KIn, VIn any, KOut comparable, VOut any](fn KVtoKV[KIn, VIn, KOut, VOut]) KeyValAll[KIn, VIn, lmap.KeyVal[KOut, VOut]] {
	return func(k KIn, v VIn) lmap.KeyVal[KOut, VOut] {
		return lmap.NewKV(fn(k, v))
	}
}

func OnVal[K comparable, VIn, VOut any](fn ValAll[VIn, VOut]) KeyValAll[K, VIn, lmap.KeyVal[K, VOut]] {
	return func(k K, v VIn) lmap.KeyVal[K, VOut] {
		return lmap.NewKV(k, fn(v))
	}
}

func OnKey[V any, KIn, KOut comparable](fn ValAll[KIn, KOut]) KeyValAll[KIn, V, lmap.KeyVal[KOut, V]] {
	return func(k KIn, v V) lmap.KeyVal[KOut, V] {
		return lmap.NewKV(fn(k), v)
	}
}

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
