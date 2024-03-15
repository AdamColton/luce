package cmpr

import (
	"golang.org/x/exp/constraints"
)

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func MinN[T constraints.Ordered](ts ...T) T {
	return compound(Min, ts)
}

func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func MaxN[T constraints.Ordered](ts ...T) T {
	return compound(Max, ts)
}

func compound[T constraints.Ordered](fn func(a, b T) T, ts []T) (t T) {
	if len(ts) == 0 {
		return
	}
	t = ts[0]
	for _, ti := range ts[1:] {
		t = fn(t, ti)
	}
	return
}

// AssertEqualizer allows a type to define an equality test.
//
// Note that when fulfilling an interface Go will coerce a pointer type to it's
// base type, but not the other way. So if AssertEqual is on the base type
// is on the base type and a pointer to that type is passed into geomtest.Equal
// it will be cast to the base type.
type AssertEqualizer interface {
	AssertEqual(to interface{}, t Tolerance) error
}
