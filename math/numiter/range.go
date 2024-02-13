// Package numiter provides ranges of numbers that can be iterated over or
// indexed as a list.
package numiter

import (
	"math"
	"reflect"

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
