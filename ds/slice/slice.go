package slice

import (
	"reflect"
	"sort"

	"github.com/adamcolton/luce/util/liter"
)

// Slice is a wrapper that provides helper methods
type Slice[T any] []T

// New is syntactic sugar to infer the type
func New[T any](s []T) Slice[T] {
	return s
}

// Make a slice with the specified length and capacity. If capacity is set to
// 0, then ln will be used for capacity as well.
func Make[T any](ln, cp int) Slice[T] {
	if cp == 0 {
		cp = ln
	}
	return make(Slice[T], ln, cp)
}

// Clone a slice.
func (s Slice[T]) Clone() Slice[T] {
	out := make([]T, len(s))
	copy(out, s)
	return out
}

// Swaps two values in the slice.
func (s Slice[T]) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Len creates a strongly typed version of builtin len for slices.
func Len[T any](s []T) int {
	return len(s)
}

// AppendNotZero will append any values from ts that are not the zero
// value for the type. Particularly useful for appending not nil values.
func (s Slice[T]) AppendNotZero(ts ...T) []T {
	for _, t := range ts {
		v := reflect.ValueOf(t)
		if v.Kind() != reflect.Invalid && !v.IsZero() {
			s = append(s, t)
		}
	}
	return s
}

// Iter returns an iter.Wrapper for the slice.
func (s Slice[T]) Iter() liter.Wrapper[T] {
	return NewIter(s)
}

// IterFactory fulfills iter.Factory.
func (s Slice[T]) IterFactory() (i liter.Iter[T], t T, done bool) {
	i = NewIter(s)
	t, done = i.Cur()
	return
}

// Remove values at given indicies by swapping them with values from the end
// and truncating the slice. Values less than zero or greater than the length
// of the list are ignored. Note that idxs is reordered so if that is a slice
// passed in and the order is important, pass in a copy.
func (s Slice[T]) Remove(idxs ...int) Slice[T] {
	sort.Sort(sort.Reverse(sort.IntSlice(idxs)))
	ln := len(s)
	prev := ln
	// Depending on variations in the implementation there are two things that
	// can make this behave in unintended ways. Duplicate values cause a double
	// swap. And it could be possible for a value near the end of the list to
	// removed, but then swapped with a value earlier in the list, reintroducing
	// it. Also, negative values are not allowed.
	//
	// To avoid both, idxs is sorted in descending order and prev tracks the
	// the last value. The "idx < prev" comparison guarentees both that there
	// are no duplicates and that idx is less than the length of the list.
	for _, idx := range idxs {
		if idx >= 0 && idx < prev {
			ln--
			s.Swap(idx, ln)
			prev = idx
		}
	}
	return s[:ln]
}

// RemoveOrdered preserves the order of the slice while removing the values
// at the given indexes.
func (s Slice[T]) RemoveOrdered(idxs ...int) Slice[T] {
	sort.Ints(idxs)
	ln := len(idxs)
	start := 0
	var pIdx int
	for {
		if start >= ln {
			return s
		}
		pIdx = idxs[start]
		start++
		if pIdx >= 0 {
			break
		}
	}
	ln = len(s)
	d := 0
	for _, idx := range idxs[start:] {
		if idx >= ln {
			break
		}
		if idx < 0 || idx == pIdx {
			continue
		}
		copy(s[pIdx-d:], s[pIdx+1:idx])
		d++
		pIdx = idx
	}
	copy(s[pIdx-d:], s[pIdx+1:ln])
	return s[:ln-d-1]
}

// Buffer is syntactic sugar to convert a Slice to a Buffer providing a set
// of methods useful for buffering operations.
func (s Slice[T]) Buffer() Buffer[T] {
	return Buffer[T](s)
}

// Pop returns the last element of the slice and the slice resized to remove
// that element. If the size of the slice is zero, the zero value for the type
// is returned.
func (s Slice[T]) Pop() (T, Slice[T]) {
	ln := len(s)
	if ln == 0 {
		var t T
		return t, s
	}
	ln--
	return s[ln], s[:ln]
}

// Shift returns the first element of the slice and the slice resized to remove
// that element. If the size of the slice is zero, the zero value for the type
// is returned.
func (s Slice[T]) Shift() (T, Slice[T]) {
	ln := len(s)
	if ln == 0 {
		var t T
		return t, s
	}
	return s[0], s[1:ln]
}

// CheckCapacity ensures that Slice s has capacity c. If not, a new slice is
// created with capacity c and the slice is copied.
func (s Slice[T]) CheckCapacity(c int) Slice[T] {
	if cap(s) >= c {
		return s
	}
	out := make(Slice[T], len(s), c)
	copy(out, s)
	return out
}

// Search wraps sort.Search
func (s Slice[T]) Search(fn func(T) bool) int {
	return sort.Search(len(s), func(idx int) bool {
		return fn(s[idx])
	})
}

// IdxCheck returns false if idx is out of the range of s.
func (s Slice[T]) IdxCheck(idx int) bool {
	return idx >= 0 && idx < len(s)
}

// Sort wraps slice.Sort. Sorts the Slice in place. The slice is also returned
// for chaining.
func (s Slice[T]) Sort(less Less[T]) Slice[T] {
	return less.Sort(s)
}

// Transform one slice to another. The transformation function's second return
// is a bool indicating if the returned value should be included in the result.
// The returned Slice is sized exactly to the output.
func Transform[In, Out any](in liter.Iter[In], fn func(In, int) (Out, bool)) Slice[Out] {
	return transform(in, fn)
}

// Transform one slice to another. The transformation function's second return
// is a bool indicating if the returned value should be included in the result.
// The returned Slice is sized exactly to the output.
func TransformSlice[In, Out any](in []In, fn func(In, int) (Out, bool)) Slice[Out] {
	return transform(NewIter(in), fn)
}

func transform[In, Out any](in liter.Iter[In], fn func(In, int) (Out, bool)) (out Slice[Out]) {
	i, done := in.Cur()
	if done {
		return
	}
	if o, include := fn(i, in.Idx()); include {
		out = transformRecurse(1, in, fn)
		out[0] = o
	} else {
		out = transformRecurse(0, in, fn)
	}
	return
}

func transformRecurse[In, Out any](size int, in liter.Iter[In], fn func(In, int) (Out, bool)) Slice[Out] {
	for i, done := in.Next(); !done; i, done = in.Next() {
		o, include := fn(i, in.Idx())
		if include {
			out := transformRecurse(size+1, in, fn)
			out[size] = o
			return out
		}
	}
	if size == 0 {
		return nil
	}
	return make([]Out, size)
}
