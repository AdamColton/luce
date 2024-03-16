package numiter

import (
	"math"

	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/math/ints"
	"golang.org/x/exp/constraints"
)

// Number is any float or integer (signed or unsigned).
type Number interface {
	constraints.Float | constraints.Integer
}

// Range with a Start, End and Step. End is exclusive.
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

// Include creates a new range making sure that end is included but that it
// is the last value in the range.
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
// of the list.List interface.
func (r *Range[T]) AtIdx(idx int) T {
	return r.Start + T(idx)*r.Step
}

// Len returns the number of steps from start to end. This fulfills a portion of
// the list.List interface.
func (r *Range[T]) Len() int {
	return int((r.End - r.Start) / r.Step)
}

func (r *Range[T]) Wrap() list.Wrapper[T] {
	return list.Wrapper[T]{r}
}

const ErrBadGrid = lerr.Str("args must be multiple of 3")

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

func IntGrid[T Number](args ...T) list.Wrapper[[]T] {
	rs := make([]list.List[T], len(args))
	for i, end := range args {
		rs[i] = IntRange(end)
	}

	return list.SliceCombinator(ints.Cross[int], rs...)
}
