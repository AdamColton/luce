package numiter

import (
	"math"

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
