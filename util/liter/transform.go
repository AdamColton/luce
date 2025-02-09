package liter

// TransformFunc converts one type to another and can be applied to an iterator.
type TransformFunc[In, Out any] func(In, int) (out Out, include bool)

// NewTransformFunc is a helper for fn without needing to explicitly specify the
// parameter types.
func NewTransformFunc[In, Out any](fn func(In, int) (out Out, include bool)) TransformFunc[In, Out] {
	return fn
}

// ForAll is a helper function for transformers. Transformers have an index
// argument and a bool return that are often not used.
func ForAll[In, Out any](fn func(in In) Out) TransformFunc[In, Out] {
	return func(in In, idx int) (Out, bool) {
		return fn(in), true
	}
}

// Transformer applies the TransformFunc to all the values of the underlying
// Iter. It includes all the values for which include is true.
type Transformer[In, Out any] struct {
	Iter[In]
	TransformFunc[In, Out]
	idx  int
	cur  Out
	done bool
}

// Factory creates a new Factory that will produce a Transformer using the
// underlying Factory and the TransformFunc.
func (fn TransformFunc[In, Out]) Factory(f Factory[In]) Factory[Out] {
	return func() (iter Iter[Out], o Out, done bool) {
		it := fn.new(f())
		return it, it.cur, it.done
	}
}

// New creates a Transformer using the give Iter with the TransformFunc.
func (fn TransformFunc[In, Out]) New(i Iter[In]) Wrapper[Out] {
	in, done := i.Cur()
	return Wrapper[Out]{fn.new(i, in, done)}
}

func (fn TransformFunc[In, Out]) new(i Iter[In], in In, done bool) *Transformer[In, Out] {
	t := &Transformer[In, Out]{
		Iter:          i,
		TransformFunc: fn,
		done:          done,
	}
	if !t.done {
		o, ok := fn(in, i.Idx())
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
func (t *Transformer[In, Out]) Next() (o Out, done bool) {
	for {
		var i In
		i, done = t.Iter.Next()
		if done {
			t.done = done
			return
		}
		var ok bool
		o, ok = t.TransformFunc(i, t.Iter.Idx())
		if ok {
			t.cur = o
			t.idx++
			return
		}
	}
}

// Cur fulfills Iter.
func (t *Transformer[In, Out]) Cur() (o Out, done bool) {
	return t.cur, t.done
}

// Done fulfills Iter.
func (t *Transformer[In, Out]) Done() bool {
	return t.done
}

// Idx fulfills Iter. The index is relative to the Transformer and not the
// underlying Iter - so it will not increment when skipping underlying values.
func (t *Transformer[In, Out]) Idx() int {
	return t.idx
}
