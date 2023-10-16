package cmpr

import "golang.org/x/exp/constraints"

// Min returns the lesser value of a or b.
func Min[T constraints.Ordered](a, b T) T {
	if b < a {
		return b
	}
	return a
}

// MinN returns the least value in ts.
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

// MaxN returns the least value in ts.
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
