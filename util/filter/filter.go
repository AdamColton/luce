package filter

// Filter is a func that tests a value of type T. Filters can be combined with
// And, Or and Not and used to select values from slices, channels and
// iterators.
type Filter[T any] func(T) bool
