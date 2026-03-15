package morph

import (
	"github.com/adamcolton/luce/util/liter"
)

// Iter applies the TransformFunc to all the values of the underlying
// Iter. It includes all the values for which include is true.
type Iter[In, Out any] struct {
	liter.Iter[In]
	Morph Val[In, Out]
	idx   int
	cur   Out
	done  bool
}

// Factory creates a new Factory that will produce a Transformer using the
// underlying Factory and the TransformFunc.
func (vt Val[In, Out]) Factory(f liter.Factory[In]) liter.Factory[Out] {
	return func() (iter liter.Iter[Out], o Out, done bool) {
		it := vt.new(f())
		return it, it.cur, it.done
	}
}

// New creates a Transformer using the give Iter with the TransformFunc.
func (vt Val[In, Out]) Iter(i liter.Iter[In]) liter.Wrapper[Out] {
	in, done := i.Cur()
	return liter.Wrapper[Out]{vt.new(i, in, done)}
}

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

// Cur fulfills Iter.
func (t *Iter[In, Out]) Cur() (o Out, done bool) {
	return t.cur, t.done
}

// Done fulfills Iter.
func (t *Iter[In, Out]) Done() bool {
	return t.done
}

// Idx fulfills Iter. The index is relative to the Transformer and not the
// underlying Iter - so it will not increment when skipping underlying values.
func (t *Iter[In, Out]) Idx() int {
	return t.idx
}
