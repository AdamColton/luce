package slice

import (
	"sort"

	"golang.org/x/exp/constraints"
)

// Less is used for sorting. Less does not inheirently mean "less than", only
// that when sorting the index should be less. For instance, when sorting from
// highest to lowest, the index of the greater value will be less that the value
// it is compared against.
type Less[T any] func(i, j T) bool

// Sort the slice using the Less function.
func (l Less[T]) Sort(s []T) []T {
	sort.Slice(s, func(i, j int) bool {
		return l(s[i], s[j])
	})
	return s
}

// LT returns an instance of Less[T] that does a less than comparison.
func LT[T constraints.Ordered]() Less[T] {
	return func(i, j T) bool {
		return i < j
	}
}

// LT returns an instance of Less[T] that does a greater than comparison.
func GT[T constraints.Ordered]() Less[T] {
	return func(i, j T) bool {
		return i > j
	}
}
