package slice

// Slice is a wrapper that provides helper methods
type Slice[T any] []T

// New is syntactic sugar to infer the type
func New[T any](s []T) Slice[T] {
	return s
}
