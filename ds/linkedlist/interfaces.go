package linkedlist

type Node[T any] interface {
	Copy() Node[T]
	Nil() bool
	Prev()
	Next()
	Get() T
	Set(payload T)
	Head()
	Tail()
	Append(payload ...T)
	Prepend(payload ...T)
}
