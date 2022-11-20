package slice

// Clone a slice.
func Clone[T any](s []T) []T {
	out := make([]T, len(s))
	copy(out, s)
	return out
}

// Swaps two values in the slice.
func Swap[T any](s []T, i, j int) {
	s[i], s[j] = s[j], s[i]
}
