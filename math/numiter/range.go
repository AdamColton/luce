// Package numiter provides ranges of numbers that can be iterated over or
// indexed as a list.
package numiter

import (
	"math"
	"reflect"

	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/math/ints"
	"github.com/adamcolton/luce/util/liter"
	"github.com/adamcolton/luce/util/reflector"
	"golang.org/x/exp/constraints"
)

// Number is any float or integer (signed or unsigned).
type Number interface {
	constraints.Float | constraints.Integer
}

// Range with a Start, End and Step. End is exclusive. A Range fulfills
// list.List (AtIdx and Len). For an integer type, Len is only correct when Step
// is positive and Start is less than End.
type Range[T Number] struct {
	Start, End, Step T
}

// NewRange creates a range from start, end and step.
func NewRange[T Number](start, end, step T) *Range[T] {
	return &Range[T]{
		Start: start,
		End:   end,
		Step:  step,
	}
}

// Index creates a Range instance from a slice.Index.
func Index(idx slice.Index) *Range[int] {
	return &Range[int]{
		Start: idx[0],
		End:   idx[1],
		Step:  1,
	}
}

// Steps creates a range between start and end with a defined number of steps.
// If includeEnd is true, the last value will be end, If it is false, the
// range will stop one step before end. The step is (end-start)/steps, which
// truncates for an integer type. steps must be at least 1, or at least 2 when
// includeEnd is true, otherwise the step is infinite or the division is by zero.
func Steps[T Number](start, end T, steps uint, includeEnd bool) *Range[T] {
	if includeEnd {
		steps--
	}
	step := (end - start) / T(steps)
	if includeEnd {
		end += step
	}
	return &Range[T]{
		Start: start,
		End:   end,
		Step:  step,
	}
}

// Include creates a new range that reaches end: the last value is the first step
// at or past end, so it is end when step divides the distance evenly and
// beyond end when it does not (Include(0, 10, 3) ends at 12).
func Include[T Number](start, end, step T) *Range[T] {
	d := float64(end - start)
	s64 := float64(step)
	steps := math.Ceil(d/s64) + 1
	end = start + T(steps*s64)
	return NewRange(start, end, step)
}

// IntRange creates a range from 0 to end with a step of 1.
func IntRange[T Number](end T) *Range[T] {
	return &Range[T]{
		End:  end,
		Step: 1,
	}
}

// AtIdx returns the value at the given step index. This fulfills a portion
// of the list.List interface. It does not check that idx is in the range.
func (r *Range[T]) AtIdx(idx int) T {
	return r.Start + T(idx)*r.Step
}

// Inv is the inverse of AtIdx. It returns the step index for t, truncating
// toward zero if t is not exactly on a step. It does not check that t is in the
// range.
func (r *Range[T]) Inv(t T) int {
	return int((t - r.Start) / r.Step)
}

// Len returns the number of steps from start to end. This fulfills a portion of
// the list.List interface.
func (r *Range[T]) Len() int {
	return divUp(r.End-r.Start, r.Step)
}

func divUp[T Number](a, b T) int {
	switch reflector.Type[T]().Kind() {
	case reflect.Float32, reflect.Float64:
		return int(math.Ceil(float64(a) / float64(b)))
	default:
		return int((a + b - 1) / b)
	}
}

// Wrap the Range to add list.Wrapper methods.
func (r *Range[T]) Wrap() list.Wrapper[T] {
	return list.Wrapper[T]{r}
}

// Iter returns an iterator over the values of the Range.
func (r *Range[T]) Iter() liter.Wrapper[T] {
	return list.NewIter(r)
}

// ErrBadGrid is the value Grid panics with when the number of args is not a
// multiple of 3.
const ErrBadGrid = lerr.Str("args must be multiple of 3")

// Grid takes args in sets of 3 as (start1, end1, step1, start2, end2, step2...)
// and returns every combination of the values of the ranges, as slices with one
// value from each range. The first range changes fastest. It panics with
// ErrBadGrid if the number of args is not a multiple of 3.
func Grid[T Number](args ...T) list.Wrapper[[]T] {
	ln := len(args)
	if ln%3 != 0 {
		panic(ErrBadGrid)
	}
	ln /= 3
	rs := make([]list.List[T], ln)
	for i := range rs {
		idx := i * 3
		rs[i] = &Range[T]{
			Start: args[idx],
			End:   args[idx+1],
			Step:  args[idx+2],
		}
	}

	return list.SliceCombinator(ints.Cross[int], rs...)
}

// NGrid is Grid with the same Range used n times: it returns every combination
// of n values of rng, as slices of length n. The first position changes
// fastest.
func NGrid[T Number](rng *Range[T], n uint) list.Wrapper[[]T] {
	rs := make([]list.List[T], n)
	for i := range n {
		rs[i] = rng
	}
	return list.SliceCombinator(ints.Cross[int], rs...)
}

// IntGrid takes the end of each range, which is from 0 to end with a step of 1,
// and returns every combination like Grid.
func IntGrid[T Number](args ...T) list.Wrapper[[]T] {
	rs := make([]list.List[T], len(args))
	for i, end := range args {
		rs[i] = IntRange(end)
	}

	return list.SliceCombinator(ints.Cross[int], rs...)
}
