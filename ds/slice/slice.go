package slice

type Slice[T any] []T

func New[T any](s []T) Slice[T] {
	return s
}
