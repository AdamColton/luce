package slice

// Slice is a generic slice type that provides helper methods.
type Slice[T any] []T

// New converts s to a Slice, inferring the type.
func New[T any](s []T) Slice[T] {
	return s
}
