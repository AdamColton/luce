package list

import (
	"reflect"

	"github.com/adamcolton/luce/util/iter"
	"github.com/adamcolton/luce/util/upgrade"
)

// Iter fulfills iter.Iter by iterating over the list.
type Iter[T any] struct {
	List[T]
	I int
}

// NewIter creates an Iter using the provided list.
func NewIter[T any](l List[T]) *Iter[T] {
	return &Iter[T]{List: l}
}

// Idx fulfills iter.Iter and provides the current index.
func (i *Iter[T]) Idx() int {
	return i.I
}

// Done fulfills iter.Iter indicating if iteration is done.
func (i *Iter[T]) Done() bool {
	return i.I >= i.Len()
}

// Cur fulfills iter.Iter returning both the current value and if iteration is
// done. If iteration is done, T will be the default value.
func (i *Iter[T]) Cur() (t T, done bool) {
	done = i.Done()
	if !done {
		t = i.List.AtIdx(i.I)
	}
	return
}

// Next fulfills iter.Iter and increments the index. It returns the current
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

// Upgrade fulfills upgrade.Upgrader allowing the underlying List to be
// upgraded.
func (i *Iter[T]) Upgrade(t reflect.Type) interface{} {
	return upgrade.Wrapped(i.List, t)
}

// IterFactory creates an iter.Factory that generates a *list.Iter backed by
// list L.
func IterFactory[T any](l List[T]) iter.Factory[T] {
	return func() (it iter.Iter[T], t T, done bool) {
		it = &Iter[T]{
			List: l,
			I:    -1,
		}
		t, done = it.Next()
		return
	}
}
