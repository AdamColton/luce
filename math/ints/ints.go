package ints

import "golang.org/x/exp/constraints"

// Dir is the direction Div rounds in when b does not divide a evenly.
type Dir byte

// The directions Div can round in.
const (
	// Down rounds toward negative infinity.
	Down Dir = iota
	// Up rounds toward positive infinity.
	Up
	// In rounds toward zero, the same as Go's integer division.
	In
	// Out rounds away from zero.
	Out
)

// Div returns a/b rounded in the direction d, whatever the signs of a and b.
// A d other than Down, Up, In or Out rounds like In. It panics if b is 0.
func Div[T constraints.Integer](a, b T, d Dir) T {
	q := a / b
	if a%b == 0 {
		return q
	}
	// Go has rounded q toward zero, which is up for a negative quotient and down
	// for a positive one.
	negative := (a < 0) != (b < 0)
	switch {
	case negative && (d == Down || d == Out):
		q--
	case !negative && (d == Up || d == Out):
		q++
	}
	return q
}

// DivDown returns a/b rounded toward negative infinity. It panics if b is 0.
func DivDown[T constraints.Integer](a, b T) T {
	return Div(a, b, Down)
}

// DivUp returns a/b rounded toward positive infinity. It panics if b is 0.
func DivUp[T constraints.Integer](a, b T) T {
	return Div(a, b, Up)
}

// DivIn returns a/b rounded toward zero, the same as Go's integer division. It
// panics if b is 0.
func DivIn[T constraints.Integer](a, b T) T {
	return Div(a, b, In)
}

// DivOut returns a/b rounded away from zero. It panics if b is 0.
func DivOut[T constraints.Integer](a, b T) T {
	return Div(a, b, Out)
}

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
