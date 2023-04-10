package list

// TransformAny takes a transfrom function without the index argument and
// converts it to a TransformFunc.
func TransformAny[In, Out any](fn func(In) Out) TransformFunc[In, Out] {
	return func(in In, idx int) (out Out) {
		return fn(in)
	}
}

// TransformFunc converts one type to another and can be applied to an iterator.
// Unlike the TransformFunc in liter and slice, this one cannot have an include
// bool because transformations on a slice must be done one-to-one.
type TransformFunc[In, Out any] func(in In, idx int) (out Out)

// Transformer applies a function to a list to transform it's values.
type Transformer[In, Out any] struct {
	List[In]
	Fn TransformFunc[In, Out]
}

// New creates Transformer using the TransformFunc and list l.
func (fn TransformFunc[In, Out]) New(l List[In]) Wrapper[Out] {
	return Transformer[In, Out]{
		List: l,
		Fn:   fn,
	}.Wrap()
}

// NewTransformer creates Transformer using TransformFunc fn and list l.
func NewTransformer[In, Out any](l List[In], fn TransformFunc[In, Out]) Wrapper[Out] {
	return fn.New(l)
}

// Note that Transformer should not have an Upgrade for the same reason as
// Reverse.

// AtIdx fulfills List by passing the value at idx in the underlying list into
// Fn.
func (t Transformer[In, Out]) AtIdx(idx int) Out {
	return t.Fn(t.List.AtIdx(idx), idx)
}

// Wrap the Transformer to add Wrapper methods.
func (t Transformer[In, Out]) Wrap() Wrapper[Out] {
	return Wrapper[Out]{t}
}
