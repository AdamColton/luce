package numiter

import (
	"math"

	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/math/ints"
	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Float | constraints.Integer
}

type Range[T Number] struct {
	Start, End, Step T
}

func NewRange[T Number](start, end, step T) *Range[T] {
	return &Range[T]{
		Start: start,
		End:   end,
		Step:  step,
	}
}

func Include[T Number](start, end, step T) *Range[T] {
	d := float64(end - start)
	s64 := float64(step)
	steps := math.Ceil(d/s64) + 1
	end = start + T(steps*s64)
	return NewRange(start, end, step)
}

func IntRange[T Number](end T) *Range[T] {
	return &Range[T]{
		End:  end,
		Step: 1,
	}
}

func (r *Range[T]) AtIdx(idx int) T {
	return r.Start + T(idx)*r.Step
}

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
