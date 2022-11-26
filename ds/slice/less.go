package slice

import "sort"

// Less is used for sorting.
type Less[T any] func(i, j T) bool

// Sort the slice using the Less function.
func (l Less[T]) Sort(s []T) {
	sort.Slice(s, func(i, j int) bool {
		return l(s[i], s[j])
	})
}
