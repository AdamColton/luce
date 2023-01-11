package list

import "github.com/adamcolton/luce/util/liter"

// Iter fulfills liter.Iter by iterating over the list.
type Iter[T any] struct {
	List[T]
	I int
}

// NewIter creates a Wrapped Iter using l.
func NewIter[T any](l List[T]) liter.Wrapper[T] {
	return liter.Wrap(&Iter[T]{
		List: l,
	})
}

// Idx fulfills liter.Iter and provides the current index.
func (i *Iter[T]) Idx() int {
	return i.I
}

// Done fulfills liter.Iter indicating if iteration is done.
func (i *Iter[T]) Done() bool {
	return i.I >= i.Len()
}

// Cur fulfills liter.Iter returning both the current value and if iteration is
// done. If iteration is done, T will be the default value.
func (i *Iter[T]) Cur() (t T, done bool) {
	done = i.Done()
	if !done {
		t = i.List.AtIdx(i.I)
	}
	return
}

// Next fulfills liter.Iter and increments the index. It returns the current
// index and a bool indicating if it's done.
func (i *Iter[T]) Next() (t T, done bool) {
	ln := i.Len()
	done = i.I >= ln
	if done {
		return
	}
	i.I++
	done = i.I >= ln
	if !done {
		t = i.AtIdx(i.I)
	}
	return
}

// Start sets the index to zero. Returns the first element and a bool indicating
// if it's done.
func (i *Iter[T]) Start() (t T, done bool) {
	i.I = -1
	return i.Next()
}

// Wrapped fulfills upgrade.Wrapper allowing the underlying List to be upgraded.
func (i *Iter[T]) Wrapped() any {
	return i.List
}
