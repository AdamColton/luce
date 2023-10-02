package filter

import "github.com/adamcolton/luce/util/iter"

type Iter[T any] struct {
	In iter.Iter[T]
	Filter[T]
	idx int
}

func (f Filter[T]) Iter(i iter.Iter[T]) iter.Wrapper[T] {
	for t, done := i.Cur(); !done && !f(t); t, done = i.Next() {
	}
	return iter.Wrap(&Iter[T]{
		In:     i,
		Filter: f,
	})
}

func (i *Iter[T]) Next() (t T, done bool) {
	for t, done = i.In.Next(); !done && !i.Filter(t); t, done = i.In.Next() {
	}
	i.idx++
	return
}

func (i *Iter[T]) Cur() (t T, done bool) {
	return i.In.Cur()
}

func (i *Iter[T]) Done() bool {
	return i.In.Done()
}

func (i *Iter[T]) Idx() int {
	return i.idx
}
