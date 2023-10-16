package ints

import "golang.org/x/exp/constraints"

// DivUp returns a/b rounding up.
func DivUp[T constraints.Integer](a, b T) T {
	out := a / b
	if a%b != 0 {
		out++
	}
	return out
}

// DivDown returns a/b rounding down. This is the Go default, but defining the
// desired behavior can be more explicit.
func DivDown[T constraints.Integer](a, b T) T {
	return a / b
}
