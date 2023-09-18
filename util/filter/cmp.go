package filter

import "golang.org/x/exp/constraints"

// EQ returns a Filter that will check if a given value is equal to 'a'.
func EQ[T constraints.Ordered](a T) Filter[T] {
	return func(b T) bool { return a == b }
}

// EQ returns a Filter that will check if a given value is greater than to 'a'.
func GT[T constraints.Ordered](a T) Filter[T] {
	return func(b T) bool { return a < b }
}

// NEQ returns a Filter that will check if a given value is not equal to 'a'.
func NEQ[T constraints.Ordered](a T) Filter[T] {
	return func(b T) bool { return a != b }
}

// LTE returns a Filter that will check if a given value is less than or equal
// to 'a'.
func LTE[T constraints.Ordered](a T) Filter[T] {
	return func(b T) bool { return a >= b }
}

// LT returns a Filter that will check if a given value is not less than to 'a'.
func LT[T constraints.Ordered](a T) Filter[T] {
	return func(b T) bool { return a > b }
}

// GTE returns a Filter that will check if a given value is greater than or
// equal to 'a'.
func GTE[T constraints.Ordered](a T) Filter[T] {
	return func(b T) bool { return a <= b }
}
