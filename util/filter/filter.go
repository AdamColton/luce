package filter

// Filter represents logic to classify a type as passing or failing.
type Filter[T any] func(T) bool
