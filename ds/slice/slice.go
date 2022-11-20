package slice

// Slice is a wrapper that provides helper methods
type Slice[T any] []T

// New is syntactic sugar to infer the type
func New[T any](s []T) Slice[T] {
	return s
}

// Clone a slice.
func (s Slice[T]) Clone() Slice[T] {
	out := make([]T, len(s))
	copy(out, s)
	return out
}

// Swaps two values in the slice.
func (s Slice[T]) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Keys returns the keys of a map as a slice
func Keys[K comparable, V any](m map[K]V) Slice[K] {
	out := make([]K, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
