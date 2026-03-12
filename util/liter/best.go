package liter

import (
	"cmp"
)

// Best finds the "best" value in the iterator using fn to convert each value to
// a comparable value then using cmpFn to compare the values. The cmpFn should
// return true is 'a' is better than 'b'.
func Best[T, C any](i Iter[T], fn func(T) C, cmpFn func(a, b C) bool) (t T, c C, idx int) {
	var best struct {
		c   C
		t   T
		idx int
	}
	best.idx = -1
	Each(i, func(idx int, t T, done *bool) {
		c := fn(t)
		if best.idx < 0 || cmpFn(c, best.c) {
			best.t, best.c, best.idx = t, c, idx
		}
	})
	return best.t, best.c, best.idx
}

// Greatest finds the value in iter that returns the greatest value from fn.
func Greatest[T any, C cmp.Ordered](i Iter[T], fn func(T) C) (t T, c C, idx int) {
	return Best(i, fn, func(a, b C) bool { return a > b })
}

// Least finds the value in iter that returns the least value from fn.
func Least[T any, C cmp.Ordered](i Iter[T], fn func(T) C) (t T, c C, idx int) {
	return Best(i, fn, func(a, b C) bool { return a < b })
}
