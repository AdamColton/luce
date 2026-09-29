package slice

import (
	"github.com/adamcolton/luce/util/upgrade"
)

// getLn returns the length of e if it fulfills Lener, possibly through a
// Wrapper, otherwise 0. It is used to size the output buffer.
func getLn(e any) int {
	if l, ok := upgrade.To[Lener](e); ok {
		return l.Len()
	}
	return 0
}

// EacherKey collects the keys from an Eacher into a Slice. If e fulfills
// Lener, it is used to set the capacity. Because Eacher does not guarantee
// order, neither does EacherKey.
func EacherKey[K, V any](e Eacher[K, V]) Slice[K] {
	out := make(Slice[K], 0, getLn(e))
	e.Each(func(k K, v V, done *bool) {
		out = append(out, k)
	})
	return out
}

// EacherVal collects the values from an Eacher into a Slice. If e fulfills
// Lener, it is used to set the capacity. Because Eacher does not guarantee
// order, neither does EacherVal.
func EacherVal[K, V any](e Eacher[K, V]) Slice[V] {
	out := make(Slice[V], 0, getLn(e))
	e.Each(func(k K, v V, done *bool) {
		out = append(out, v)
	})
	return out
}
