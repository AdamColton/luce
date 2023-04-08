package list

// List behaves like an array or slice - it has index values and a length.
type List[T any] interface {
	// AtIdx returns the value at idx.
	AtIdx(idx int) T
	// Len returns the number of values in the List.
	Len() int
}

// Slicer allows Lists to provide efficient methods for generating slices. If a
// List fulfills Slicer, Wrapper.Slice will use that to generate a slice instead
// of iterating over the List. It is the same as slice.Slicer.
type Slicer[T any] interface {
	Slice(buf []T) []T
}
