package filter

// Filter provides tools to filter ints and compose filters
type Filter[T any] func(T) bool
