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
