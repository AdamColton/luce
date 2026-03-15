package morph

import (
	"github.com/adamcolton/luce/util/liter"
)

// Iter applies the Val to all the values of the underlying Iter. It includes
// all the values for which include is true.
type Iter[In, Out any] struct {
	liter.Iter[In]
	Morph Val[In, Out]
	idx   int
	cur   Out
	done  bool
}

// Factory creates a new Factory that will produce an Iter using the underlying
// Factory and the Val.
func (vt Val[In, Out]) Factory(f liter.Factory[In]) liter.Factory[Out] {
	return func() (iter liter.Iter[Out], o Out, done bool) {
		it := vt.new(f())
		return it, it.cur, it.done
	}
}

// Iter creates an Iter using the given iterator with the Val. It starts from the
// given iterator's current value, it does not reset it.
func (vt Val[In, Out]) Iter(i liter.Iter[In]) liter.Wrapper[Out] {
	in, done := i.Cur()
	return liter.Wrapper[Out]{vt.new(i, in, done)}
}

// Iter creates an Iter using the given iterator with the ValAll, which includes
// every value. It starts from the given iterator's current value.
func (vt ValAll[In, Out]) Iter(i liter.Iter[In]) liter.Wrapper[Out] {
	return vt.ToV().Iter(i)
}

func (vt Val[In, Out]) new(i liter.Iter[In], in In, done bool) *Iter[In, Out] {
	t := &Iter[In, Out]{
		Iter:  i,
		Morph: vt,
		done:  done,
	}
	if !t.done {
		o, ok := vt(in)
		if ok {
			t.cur = o
		} else {
			t.idx = -1
			t.Next()
		}
	}
	return t
}

// Next fulfills Iter, it returns the transformation of the next included
// value from the underlying iterator.
func (t *Iter[In, Out]) Next() (o Out, done bool) {
	for {
		var i In
		i, done = t.Iter.Next()
		if done {
			t.done = done
			return
		}
		var ok bool
		o, ok = t.Morph(i)
		if ok {
			t.cur = o
			t.idx++
			return
		}
	}
}

// Cur fulfills Iter. If the Iter is done, it returns the zero value.
func (t *Iter[In, Out]) Cur() (o Out, done bool) {
	if t.done {
		return o, true
	}
	return t.cur, false
}

// Done fulfills Iter.
func (t *Iter[In, Out]) Done() bool {
	return t.done
}

// Idx fulfills Iter. The index is relative to the Iter and not the underlying
// Iter - so it will not increment when skipping underlying values.
func (t *Iter[In, Out]) Idx() int {
	return t.idx
}
