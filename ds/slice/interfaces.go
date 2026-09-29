package slice

// Slicer allows an interface to show that it has an efficient way to convert
// itself to a slice.
type Slicer[T any] interface {
	Slice(buf []T) []T
}

// Lener allows an interface to show that it knows its length.
type Lener interface {
	Len() int
}

// EachFn is the callback used by Eacher. It is called once for each key/value
// pair. Setting done to true stops the iteration.
type EachFn[K, V any] = func(k K, v V, done *bool)

// Eacher is fulfilled by any type that can call a function on each of its
// key/value pairs. No order is guaranteed.
type Eacher[K, V any] interface {
	Each(EachFn[K, V])
}
