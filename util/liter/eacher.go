package liter

// Eacher turns a function that calls an EachFn for every value into something
// with Each and Len methods, so it can be used where a collection with those
// methods is expected.
type Eacher[V any] struct {
	// Func calls the EachFn it is given for every value, and stops if done is
	// set.
	Func func(inner EachFn[V])
	// L is the number of values. It is returned by Len.
	L int
}

// Each calls Func with fn.
func (e Eacher[V]) Each(fn EachFn[V]) {
	e.Func(fn)
}

// Len returns L.
func (e Eacher[V]) Len() int {
	return e.L
}
