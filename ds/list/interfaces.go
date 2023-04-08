package list

// List behaves like an array or slice - it has index values and a length.
type List[T any] interface {
	AtIdx(idx int) T
	Len() int
}

// Slicer allows Lists to provide efficient methods for generating slices.
// If a List fulfills Slicer the ToSlice function will use that to generate a
// slice instead of iterating over the List.
type Slicer[T any] interface {
	Slice(buf []T) []T
}
