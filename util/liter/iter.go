// Package liter provides an iterator. To correctly implement an iterator,
// it should be initilized in a valid state so that this for loop would
// visit all the values:
//
//	for t,done := i.Cur(); !done; t,done = i.Next(){...}
package liter

// Iter interface allows for a standard set of tools for iterating over a
// collection.
type Iter[T any] interface {
	Next() (t T, done bool)
	Cur() (t T, done bool)
	Done() bool
	Idx() int
}

// Starter is an optional interface that Iter can implement to return to the
// start of the iteration.
type Starter[T any] interface {
	Start() (t T, done bool)
}
