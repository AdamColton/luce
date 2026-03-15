package slice

import "cmp"

type Compare[T any] func(a, b T) int

func NewCompare[T any](fn func(a, b T) int) Compare[T] {
	return fn
}

func (c Compare[T]) NewOrdered(s []T) Ordered[T] {
	return Ordered[T]{
		Slice:   s,
		Compare: c,
	}.Sort()
}

type Ordered[T any] struct {
	Slice[T]
	Compare func(a, b T) int
}

func (s Slice[T]) Ordered(compare func(a, b T) int) Ordered[T] {
	return Ordered[T]{
		Slice:   s,
		Compare: compare,
	}
}

func NewOrdered[T cmp.Ordered](s []T) Ordered[T] {
	return NewCompare(cmp.Compare[T]).NewOrdered(s)
}

func (c Ordered[T]) Contains(find T) bool {
	_, found := c.Find(find)
	return found
}

func (c Ordered[T]) Find(find T) (idx int, found bool) {
	return c.Slice.Find(func(t T) int {
		return c.Compare(find, t)
	})
}

func (c Ordered[T]) Sort() Ordered[T] {
	c.Slice.Sort(func(i, j T) bool {
		return c.Compare(i, j) != 1
	})
	return c
}
