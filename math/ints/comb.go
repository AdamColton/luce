package ints

import (
	"golang.org/x/exp/constraints"
)

// Combinator represents a list of combinations that can be returned by
// index.
type Combinator[I constraints.Integer] func(idx I) []I

// CombinatorFactory takes a slice of integers and produces a combinator for
// them.
type CombinatorFactory[I constraints.Integer] func(lns ...I) (c Combinator[I], ln I)

// Pair loops through the provided indexes until they all end simultaneously.
// If the lengths are all the same, this will be a single cycle.
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

// Cross produces every combination lengths.
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
