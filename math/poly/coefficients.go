package poly

import (
	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/ds/slice"
)

// Coefficients is the list of coefficients of a polynomial. The index is the
// power of x that the coefficient belongs to, so the list 1, 2, 3 is 1+2x+3x^2.
// Indexes past the end have a coefficient of 0.
type Coefficients = list.List[float64]

// Buf creates a Slice with c capacity and a value of 1, reusing buf if it has
// the capacity. This is useful when taking the product of several polynomials.
func Buf(c int, buf []float64) Slice {
	return Slice(append(slice.NewBuffer(buf).Empty(c), 1))
}

// Slice is Coefficients held in a []float64.
type Slice []float64

// AtIdx is the coefficient at idx. If the idx is greater than the length of the
// polynomial, then a 0 is returned. It panics if idx is negative.
func (s Slice) AtIdx(idx int) float64 {
	if idx >= s.Len() {
		return 0
	}
	return s[idx]
}

// Len is the length of the slice.
func (s Slice) Len() int {
	return len(s)
}

// BufLen returns a Slice of length ln. If buf has the capacity it is used and
// the values in it are left as they were, otherwise a new Slice is made.
func BufLen(buf []float64, ln int) Slice {
	return Slice(slice.NewBuffer(buf).Slice(ln))
}

// BufSplit returns an empty Slice with capacity ln, and the rest of buf to be
// used for something else. If buf does not have the capacity ln, a new Slice is
// made and buf is returned unchanged.
func BufSplit(buf []float64, ln int) (Slice, []float64) {
	s, f := slice.NewBuffer(buf).Split(ln)
	return Slice(s), f
}

// BufEmpty returns a Slice of length 0 and capacity ln. If buf has the capacity
// it is used, otherwise a new Slice is made.
func BufEmpty(buf []float64, ln int) Slice {
	return Slice(slice.NewBuffer(buf).Empty(ln))
}
