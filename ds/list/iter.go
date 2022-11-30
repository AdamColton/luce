package list

// TEST

type Iter[T any] struct {
	List[T]
	Cur int
}

func NewIter[T any](l List[T]) *Iter[T] {
	return &Iter[T]{List: l}
}

func (i *Iter[T]) Idx() int {
	return i.Cur
}

func (i *Iter[T]) Done() bool {
	return i.Cur < i.Len()
}

func (i *Iter[T]) Next() (T, bool) {
	var t T
	ln := i.Len()
	done := i.Cur >= ln
	if !done {
		i.Cur++
		done = i.Cur >= ln
		if !done {
			t = i.AtIdx(i.Cur)
		}
	}
	return t, done
}

func (i *Iter[T]) Start() (T, bool) {
	i.Cur = -1
	return i.Next()
}

func (i *Iter[T]) Do(fn func(idx int, t T) bool) {
	for t, done := i.Start(); !done; t, done = i.Next() {
		if fn(i.Cur, t) {
			return
		}
	}
}
