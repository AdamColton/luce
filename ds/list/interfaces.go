package list

// List behaves like an array or slice - it has index values and a length.
type List[T any] interface {
	AtIdx(idx int) T
	Len() int
}
