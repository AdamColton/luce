package ints

import "golang.org/x/exp/constraints"

// DivUp returns a/b rounding up. If a and b have different signs and b does not
// divide a evenly, the result is one too large.
func DivUp[T constraints.Integer](a, b T) T {
	out := a / b
	if a%b != 0 {
		out++
	}
	return out
}

// DivDown returns a/b using Go's integer division, which rounds towards zero.
// For a and b with the same sign that is rounding down. Defining the desired
// behavior can be more explicit.
func DivDown[T constraints.Integer](a, b T) T {
	return a / b
}

// == projects.Code.luce.ints ==
// [ ] DivRound(a, b T) T

// Mod provides a version of modulus consistent with most other languages and
// calculators.
//
// In Go, mod (%) will return a negative if either a or b is negative. In most
// other languages and calculators the sign will always match b. Mod panics if b
// is 0.
func Mod[T constraints.Integer](a, b T) T {
	if a < 0 {
		m := (b - (-a % b)) % b
		return m
	}
	m := a % b
	if m > 0 && b < 0 {
		m += b
	}
	return m
}
