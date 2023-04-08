package lset

// Multi combines multiple sets and treats them as a single set
type Multi[T comparable] []*Set[T]

// NewMulti is a helper that infers type when creating a Multi
func NewMulti[T comparable](ts ...*Set[T]) Multi[T] {
	return ts
}
