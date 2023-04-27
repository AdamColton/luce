package slice

// Buffer is used to provide a slice for re-use avoiding excessive allocation.
type Buffer[T any] []T
