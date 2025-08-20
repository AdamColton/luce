package flow

import "golang.org/x/exp/constraints"

// BitFlag supports bit flag operations. The Flag can be a single bit or several
// bits.
type BitFlag[T constraints.Integer] struct {
	Flag T
}

// NewFlag creates a BitFlag with f as the flag.
func NewFlag[T constraints.Integer](f T) BitFlag[T] {
	return BitFlag[T]{Flag: f}
}

// Check returns true if the flag is set on f. If the flag has several bits, all
// of them must be set.
func (bf BitFlag[T]) Check(f T) bool {
	return bf.Flag&f == bf.Flag
}

// Set the flag on f.
func (bf BitFlag[T]) Set(f *T) {
	*f = (*f) | bf.Flag
}

// Clear the flag on f.
func (bf BitFlag[T]) Clear(f *T) {
	*f = (*f) & (^bf.Flag)
}

// OrBitFlag is a set of flags where any one of them is enough. A flag with
// several bits needs all of its bits.
type OrBitFlag[T constraints.Integer] []T

// NewOrFlag creates an OrBitFlag from flags.
func NewOrFlag[T constraints.Integer](flags ...T) OrBitFlag[T] {
	return flags
}

// Check returns true if any of the flags is set on f. It is false for an empty
// OrBitFlag, and always true if one of the flags is 0.
func (obf OrBitFlag[T]) Check(f T) bool {
	for _, bf := range obf {
		if bf&f == bf {
			return true
		}
	}
	return false
}
