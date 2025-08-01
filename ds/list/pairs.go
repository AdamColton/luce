package list

import "github.com/adamcolton/luce/math/ints"

// Pairs fulfills the List interface. Each call to AtIdx returns the value at
// that index in the underlying list and the next index as an [2]T. AtIdx uses
// modulus, so it will work for any values. The value of AtIdx for the length of
// the underlying List will return a pair containing the first and last value in
// the underlying list. For this reason, Loop actually only effects the Len().
// If Loop is true Len() returns the same value as the underlying list and if
// it is false, it returns a value one less.
type Pairs[T any] struct {
	List[T]
	Loop bool
}

// NewPairs creates an instance of *Pairs from the list. If loop is true, then
// the Len will be the same as the underlying List, if false, it will be one
// less.
func NewPairs[T any](l List[T], loop bool) *Pairs[T] {
	return &Pairs[T]{
		List: l,
		Loop: loop,
	}
}

// AtIdx returns a pair of values corresponding to idx and idx+1 of the
// underlying list. It uses modulus so that a valid pair is returned for any
// idx and so that p.AtIdx(p.List.Len()) returns a pair containing the first
// and last values in the underlying list.
func (p *Pairs[T]) AtIdx(idx int) [2]T {
	ln := p.List.Len()
	idx = ints.Mod(idx, ln)
	return [2]T{p.List.AtIdx(idx), p.List.AtIdx((idx + 1) % ln)}
}

// Len returns the number of pairs in the list. If loop is true, this is the
// same as the underlying list, if it is false, it's one less.
func (p *Pairs[T]) Len() int {
	ln := p.List.Len()
	if !p.Loop {
		ln--
	}
	return ln
}

// Wrap the Pairs to add Wrapper methods.
func (p *Pairs[T]) Wrap() Wrapper[[2]T] {
	return Wrapper[[2]T]{p}
}
