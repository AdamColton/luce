package list

type Transformer[In, Out any] struct {
	List[In]
	Fn func(In) Out
}

// Note that Transformer should not have an Upgrade for the same reason as
// Reverse.

// AtIdx fulfills List by passing the value at idx in the underlying list into
// Fn.
func (t Transformer[In, Out]) AtIdx(idx int) Out {
	return t.Fn(t.List.AtIdx(idx))
}

// Wrap the Transformer to add Wrapper methods.
func (t Transformer[In, Out]) Wrap() Wrapper[Out] {
	return Wrapper[Out]{t}
}
