package ints

import "golang.org/x/exp/constraints"

func DivUp[T constraints.Integer](a, b T) T {
	out := a / b
	if a%b != 0 {
		out++
	}
	return out
}

func DivDown[T constraints.Integer](a, b T) T {
	return a / b
}
