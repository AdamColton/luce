package flow

import "golang.org/x/exp/constraints"

// BitFlag supports bit flag operations.
type BitFlag[T constraints.Integer] struct {
	Flag T
}

// NewFlag creates a BitFlag
func NewFlag[T constraints.Integer](f T) BitFlag[T] {
	return BitFlag[T]{Flag: f}
}

// Check if bf is set on f.
func (bf BitFlag[T]) Check(f T) bool {
	return bf.Flag&f == bf.Flag
}

// Set the bit flag on f
func (bf BitFlag[T]) Set(f *T) {
	*f = (*f) | bf.Flag
}

// Clear the bit flag on f
func (bf BitFlag[T]) Clear(f *T) {
	*f = (*f) & (^bf.Flag)
}

type OrBitFlag[T constraints.Integer] []T

func NewOrFlag[T constraints.Integer](flags ...T) OrBitFlag[T] {
	return flags
}

func (obf OrBitFlag[T]) Check(f T) bool {
	for _, bf := range obf {
		if bf&f == bf {
			return true
		}
	}
	return false
}
