package filter

// Filter represents boolean logic on a Type.
type Filter[T any] func(T) bool
