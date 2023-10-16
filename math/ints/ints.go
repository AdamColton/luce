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
