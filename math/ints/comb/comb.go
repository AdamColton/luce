package comb

import (
	"github.com/adamcolton/luce/math/ints"
	"golang.org/x/exp/constraints"
)

type Combinator[I constraints.Integer] func(idx I) (a, b I)
type CombinatorFactory[I constraints.Integer] func(lnA, lnB I) (c Combinator[I], ln I)

func Pair[I constraints.Integer](lnA, lnB I) (Combinator[I], I) {
	return func(idx I) (a I, b I) {
		return ints.Mod(idx, lnA), ints.Mod(idx, lnB)
	}, ints.LCM(lnA, lnB)
}

func Cross[I constraints.Integer](lnA, lnB I) (Combinator[I], I) {
	ln := lnA * lnB
	return func(idx I) (a I, b I) {
		idx = ints.Mod(idx, ln)
		a = ints.Mod(idx, lnA)
		b = idx / lnA
		return
	}, ln
}
