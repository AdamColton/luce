package cmpr

import "golang.org/x/exp/constraints"

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
