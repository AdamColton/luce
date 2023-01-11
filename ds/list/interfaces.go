package list

// List behaves like an array or slice - it has index values and a length.
type List[T any] interface {
	// AtIdx returns the value at idx.
	AtIdx(idx int) T
	// Len returns the number of values in the List.
	Len() int
}
