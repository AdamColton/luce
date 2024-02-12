package slice

import (
	"reflect"
	"sort"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/math/cmpr"
	"github.com/adamcolton/luce/math/ints"
	"github.com/adamcolton/luce/util/liter"
)

// Slice is a generic slice type that provides helper methods.
type Slice[T any] []T

// New converts s to a Slice, inferring the type.
func New[T any](s []T) Slice[T] {
	return s
}

// Vals creates a Slice from the arguments. This is syntactic sugar for
// creating a Slice without specifying the type.
func Vals[T any](vals ...T) Slice[T] {
	return vals
}

// Make creates a Slice with the specified length and capacity. If cp is 0, ln
// is used for the capacity as well.
func Make[T any](ln, cp int) Slice[T] {
	if cp == 0 {
		cp = ln
	}
	return make(Slice[T], ln, cp)
}

// NewCap makes an empty Slice with capacity c.
func NewCap[T any](c int) Slice[T] {
	return make(Slice[T], 0, c)
}

// Clone returns a copy of the slice. The capacity can be set with cp. If cp is less than the length
// of s, that length will be used as the capacity. If cp is less than zero,
// then the length of s will be used.
func (s Slice[T]) Clone(cp int) Slice[T] {
	ln := len(s)
	if cp < 0 {
		cp = ln
	} else {
		cp = cmpr.Max(cp, ln)
	}
	out := make([]T, ln, cp)
	copy(out, s)
	return out
}

// Swap swaps two values in the slice. The indexes are relative (see Idx) so a value
// of -1 refers to the last element. It will panic if either index is out of
// range.
func (s Slice[T]) Swap(i, j int) {
	i, _ = s.Idx(i)
	j, _ = s.Idx(j)
	s[i], s[j] = s[j], s[i]
}

// Len is a strongly typed version of the builtin len for slices. It is useful
// when a func value is needed.
func Len[T any](s []T) int {
	return len(s)
}

// AppendNotZero appends any values from ts that are not the zero value for the
// type, and returns the result. Particularly useful for appending not nil
// values. A nil interface value counts as zero.
func (s Slice[T]) AppendNotZero(ts ...T) []T {
	for _, t := range ts {
		v := reflect.ValueOf(t)
		if v.Kind() != reflect.Invalid && !v.IsZero() {
			s = append(s, t)
		}
	}
	return s
}

// Iter returns a liter.Wrapper for the slice.
func (s Slice[T]) Iter() liter.Wrapper[T] {
	return NewIter(s)
}

// IterFactory fulfills liter.Factory. It returns a new iterator over the slice
// and its first value.
func (s Slice[T]) IterFactory() (i liter.Iter[T], t T, done bool) {
	i = NewIter(s)
	t, done = i.Cur()
	return
}

// Remove values at given indices by swapping them with values from the end
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
	// last value. The "idx < prev" comparison guarantees both that there
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

// RemoveOrdered preserves the order of the slice while removing the values at
// the given indexes. Negative and duplicate indexes are ignored, as are indexes
// past the end of the slice, unless the smallest non-negative index is out of
// range, in which case it panics. Note that idxs is sorted in place, so if that
// is a slice passed in and the order is important, pass in a copy.
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

// Buffer is syntactic sugar to convert a Slice to a Buffer providing a set of
// methods useful for buffering operations.
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

// CheckCapacity ensures that Slice s has a capacity of at least check. If not,
// a new slice is created with capacity request and the slice is copied.
//
// If request is 0 (or any value less than check) then check is used. Request
// can be used to avoid additional allocations. For instance, doubling the
// capacity if there is insufficient space (as the underlying slice does in Go).
func (s Slice[T]) CheckCapacity(check, request int) Slice[T] {
	if cap(s) >= check {
		return s
	}
	request = cmpr.Max(check, request)
	out := make(Slice[T], len(s), request)
	copy(out, s)
	return out
}

// Search wraps sort.Search. The slice must be sorted so that fn returns false
// for a prefix of the slice and true for the rest. Search returns the first idx
// where fn is true, or len(s) if there is none.
func (s Slice[T]) Search(fn func(T) bool) int {
	return sort.Search(len(s), func(idx int) bool {
		return fn(s[idx])
	})
}

// Find wraps sort.Find. The slice must be sorted so that compare returns
// values > 0 for a prefix of the slice, then 0, then < 0. Find returns the
// first idx where compare(s[idx]) <= 0 and found is true if compare returns 0
// at that idx. If there is no such element, idx is len(s).
func (s Slice[T]) Find(compare func(T) int) (idx int, found bool) {
	fn := func(idx int) int {
		return compare(s[idx])
	}
	return sort.Find(len(s), fn)
}

// IdxCheck returns false if idx is out of the range of s.
func (s Slice[T]) IdxCheck(idx int) bool {
	return idx >= 0 && idx < len(s)
}

// Idx provides a relative index to the slice. So a value of -1 will return
// the last index. The bool indicates if the index is in range.
func (s Slice[T]) Idx(idx int) (int, bool) {
	return ints.Idx(idx, len(s))
}

// Sort wraps Less.Sort. Sorts the Slice in place. The slice is also returned
// for chaining.
func (s Slice[T]) Sort(less Less[T]) Slice[T] {
	return less.Sort(s)
}

// Reverse a slice in place.
func (s Slice[T]) Reverse() {
	ln := len(s)
	end := ln / 2
	ln--
	for i := 0; i < end; i++ {
		s.Swap(i, ln-i)
	}
}

// ErrRng is the error AtIdx panics with when the index is out of range.
const ErrRng = lerr.Str("index out of range")

// AtIdx returns the value at idx. Fulfills list.List. Uses a relative index so
// a value of -1 will return the last index. If idx is outside the range, it
// will panic.
func (s Slice[T]) AtIdx(idx int) T {
	idx = lerr.OK(s.Idx(idx))(ErrRng)
	return s[idx]
}

// Len returns the length of the slice. Fulfills list.List.
func (s Slice[T]) Len() int {
	return len(s)
}

// AppendIf appends v to s if cond is true, otherwise s is returned unchanged.
// This is useful for conditionally adding values while building a slice in a
// single expression.
func (s Slice[T]) AppendIf(cond bool, v ...T) Slice[T] {
	if cond {
		return append(s, v...)
	}
	return s
}

// Split s at idx returning s[:idx] and s[idx:]. Both share the same backing
// array as s. It will panic if idx is outside the range [0, len(s)].
func (s Slice[T]) Split(idx int) (Slice[T], Slice[T]) {
	return s[:idx], s[idx:]
}
