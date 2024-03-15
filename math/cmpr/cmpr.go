package cmpr

import (
	"golang.org/x/exp/constraints"
)

// Min returns the lesser value of a or b.
func Min[T constraints.Ordered](a, b T) T {
	if b < a {
		return b
	}
	return a
}

// MinN returns the least value in ts, or the zero value if ts is empty.
func MinN[T constraints.Ordered](ts ...T) T {
	return compound(Min, ts)
}

// Max returns the greater value of a or b.
func Max[T constraints.Ordered](a, b T) T {
	if b > a {
		return b
	}
	return a
}

// MaxN returns the greatest value in ts, or the zero value if ts is empty.
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

// AssertEqualizer allows a type to define an equality test. AssertEqual returns
// an error if the value is not equal to the argument to within the Tolerance.
//
// Note that a pointer type fulfills an interface that its base type fulfills,
// but not the other way. So if AssertEqual is on the base type, both the base
// type and a pointer to it are AssertEqualizers.
type AssertEqualizer interface {
	AssertEqual(to interface{}, t Tolerance) error
}
