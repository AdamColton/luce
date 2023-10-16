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

func Mod[T constraints.Integer](a, b T) T {
	if a < 0 {
		return (b - (-a % b)) % b
	}
	return a % b
}
