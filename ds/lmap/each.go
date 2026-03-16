package lmap

import "github.com/adamcolton/luce/util/upgrade"

type Each[K comparable, V any] struct {
	Func func(EachFunc[K, V])
	L    int
}

func (e Each[K, V]) Each(fn EachFunc[K, V]) {
	e.Func(fn)
}

func (e Each[K, V]) Len() int {
	return e.L
}

type IdxEacher[K comparable, V any] interface {
	Each(EachFunc[int, KeyVal[K, V]])
}

type Lener interface {
	Len() int
}

func FromEach[K comparable, V any](fn func(EachFunc[int, KeyVal[K, V]])) Wrapper[K, V] {
	m := make(map[K]V)
	fn(func(idx int, kv KeyVal[K, V], done *bool) {
		m[kv.Key] = kv.Val
	})
	return Wrap(Map[K, V](m))
}

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

func (w Wrapper[K, V]) AppendEacher(e IdxEacher[K, V]) Wrapper[K, V] {
	return w.AppendEach(e.Each)
}

func (w Wrapper[K, V]) AppendEach(fn func(EachFunc[int, KeyVal[K, V]])) Wrapper[K, V] {
	fn(func(idx int, kv KeyVal[K, V], done *bool) {
		w.Set(kv.Key, kv.Val)
	})
	return w
}
