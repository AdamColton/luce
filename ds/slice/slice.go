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

// Keys returns the keys of a map as a slice
func Keys[K comparable, V any](m map[K]V) []K {
	out := make([]K, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// Vals returns the values of a map as a slice.
func Vals[K comparable, V any](m map[K]V) []V {
	out := make([]V, 0, len(m))
	for _, v := range m {
		out = append(out, v)
	}
	return out
}

// Unique returns a slice with all the unique elements of the slice passed in.
func Unique[T comparable](s []T) []T {
	set := make(map[T]struct{})
	for _, t := range s {
		set[t] = struct{}{}
	}
	return Keys(set)
}
