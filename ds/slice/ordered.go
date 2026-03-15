package slice

import "cmp"

// Compare is a comparison function. It returns a negative number if a is less
// than b, zero if they are equal and a positive number if a is greater. Ordered
// expects exactly -1, 0 and 1, as returned by cmp.Compare, and may sort
// incorrectly with other values.
type Compare[T any] func(a, b T) int

// NewCompare is syntactic sugar to infer the type of a Compare.
func NewCompare[T any](fn func(a, b T) int) Compare[T] {
	return fn
}

// NewOrdered sorts s in place using the Compare and returns it as an Ordered.
func (c Compare[T]) NewOrdered(s []T) Ordered[T] {
	return Ordered[T]{
		Slice:   s,
		Compare: c,
	}.Sort()
}

// Ordered is a Slice that is sorted according to its Compare, which makes
// searching it efficient. Nothing checks that the Slice is still sorted after
// it is modified; Sort restores the order.
type Ordered[T any] struct {
	Slice[T]
	Compare func(a, b T) int
}

// Ordered pairs the Slice with a Compare. The Slice is used as is; call Sort if
// it may not be sorted yet.
func (s Slice[T]) Ordered(compare func(a, b T) int) Ordered[T] {
	return Ordered[T]{
		Slice:   s,
		Compare: compare,
	}
}

// NewOrdered sorts s in place and returns it as an Ordered that uses
// cmp.Compare.
func NewOrdered[T cmp.Ordered](s []T) Ordered[T] {
	return NewCompare(cmp.Compare[T]).NewOrdered(s)
}

// Contains returns true if find is in the Slice.
func (c Ordered[T]) Contains(find T) bool {
	_, found := c.Find(find)
	return found
}

// Find returns the index of find and true. If find is not in the Slice it
// returns the index where find would be inserted and false. If there are
// duplicates, it returns the first.
func (c Ordered[T]) Find(find T) (idx int, found bool) {
	return c.Slice.Find(func(t T) int {
		return c.Compare(find, t)
	})
}

// Sort the Slice in place using Compare. The sort is not stable. The Ordered is
// returned for chaining.
func (c Ordered[T]) Sort() Ordered[T] {
	c.Slice.Sort(func(i, j T) bool {
		return c.Compare(i, j) != 1
	})
	return c
}
