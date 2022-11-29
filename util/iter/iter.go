package iter

// Iter interface allows for a standard set of tools for iterating over a
// collection.
type Iter[T any] interface {
	Next() (t T, done bool)
	Cur() (t T, done bool)
	Done() bool
	Idx() int
}
