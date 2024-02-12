package poly

import (
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/math/cmpr"
)

// Poly is a 1D polynomial. The index corresponds power of X.
type Poly struct {
	Coefficients
}

// New 1D polynomial with the given coefficients.
func New(cs ...float64) Poly {
	ln := len(cs)
	if ln == 0 {
		return Poly{Empty{}}
	}
	if cs[ln-1] == 0 {
		return New(cs[:ln-1]...)
	}
	if ln == 1 {
		return Poly{D0(cs[0])}
	}
	if ln == 2 && cs[1] == 1 {
		return Poly{D1(cs[0])}
	}
	return Poly{Slice(cs)}
}

// Copy a Polynomial into a buffer.
func (p Poly) Copy(buf []float64) Poly {
	out := BufLen(buf, p.Len())
	for i := range out {
		out[i] = p.AtIdx(i)
	}
	return Poly{out}
}

// Buf tries to get the Coefficients as a []float64. This is intended for
// recycling buffers.
func (p Poly) Buf() []float64 {
	buf, _ := p.Coefficients.(Slice)
	return buf
}

// F computes the value of p(x).
func (p Poly) F(x float64) float64 {
	idx := p.Len() - 1
	s := 0.0
	for ; idx >= 0; idx-- {
		s = p.AtIdx(idx) + s*x
	}
	return s
}

// AssertEqual allows Polynomials to be compared. This fulfills
// cmprtest.AssertEqualizer.
func (p Poly) AssertEqual(to any, t cmpr.Tolerance) error {
	if err := lerr.NewTypeMismatch(p, to); err != nil {
		return err
	}
	p2 := to.(Poly)

	ln := p.Len()
	if ln2 := p2.Len(); ln2 > ln {
		ln = ln2
	}
	return lerr.NewSliceErrs(ln, -1, func(i int) error {
		c0, c1 := p.AtIdx(i), p2.AtIdx(i)
		return lerr.NewNotEqual(c0 == c1, c0, c1)
	})
}
