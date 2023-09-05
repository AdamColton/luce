package slice

// Slicer allows an interface to show that it has an efficient way to convert
// itself to a slice.
type Slicer[T any] interface {
	Slice(buf []T) []T
}

// Lener allows an interface to show that it knows it's length.
type Lener interface {
	Len() int
}
