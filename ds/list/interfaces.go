package list

type List[T any] interface {
	AtIdx(idx int) T
	Len() int
}

type Slicer[T any] interface {
	Slice(buf []T) []T
}
