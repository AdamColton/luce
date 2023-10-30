package ints

import (
	"golang.org/x/exp/constraints"
)

type Combinator[I constraints.Integer] func(idx I) []I
type CombinatorFactory[I constraints.Integer] func(lns ...I) (c Combinator[I], ln I)

func Pair[I constraints.Integer](lns ...I) (Combinator[I], I) {
	// Todo: great! now this needs a buffer
	lln := len(lns)
	ln := LCMN(lns...)
	return func(idx I) []I {
		out := make([]I, lln)
		for i, ln := range lns {
			out[i] = Mod(idx, ln)
		}
		return out
	}, ln
}

func Cross[I constraints.Integer](lns ...I) (Combinator[I], I) {
	lln := len(lns)
	ln := Prod(lns...)

	return func(idx I) []I {
		idx = Mod(idx, ln)
		out := make([]I, lln)
		var s I = 1
		for i, ln := range lns {
			ii := idx / s
			s *= ln
			out[i] = Mod(ii, ln)
		}
		return out
	}, ln
}
