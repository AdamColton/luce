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

// Divide creates a new polynomial by dividing p by (x-n). The float64 returned
// is the remainder. If (x-n) is a root of p this value will be 0. Dividing a
// constant gives an empty polynomial and the constant as the remainder, and
// dividing an empty polynomial gives an empty polynomial and a remainder of 0.
func (p Poly) Divide(n float64, buf []float64) (Poly, float64) {
	ln := p.Len() - 1
	if ln < 0 {
		return Poly{Empty{}}, 0
	}
	out := BufLen(buf, ln)
	r := p.AtIdx(ln)
	for i := ln - 1; i >= 0; i-- {
		out[i], r = r, p.AtIdx(i)+r*n
	}
	return Poly{out}, r
}

// Add p and p2 using the Sum coefficients.
func (p Poly) Add(p2 Poly) Poly {
	return Poly{Sum{p, p2}}
}

// Scale will return an instace of the Scale Coefficient wrapper.
func (p Poly) Scale(s float64) Poly {
	return Poly{Scale{
		By:           s,
		Coefficients: p,
	}}
}

// Multiply two polynomails. Note that it is not safe to reuse either input as
// the buffer.
func (p Poly) Multiply(p2 Poly) Poly {
	return Poly{Product{p, p2}}
}

// MultSwap does a multiply and swap. It is used for effiency when doing
// consecutive multiplications. It is equivalent to:
//
// p = p.Multiply(p2)
//
// but it swaps the slice backing p with the buf after the multiplicaiton. It
// will generally be used like this:
//
// buf = p.MultSwap(p2, buf)
//
// Unlike the other methods, it has a pointer receiver, because it changes p.
// The slice that was backing p is returned, so it can be used as the next buf.
func (p *Poly) MultSwap(p2 Poly, buf []float64) []float64 {
	prod := p.Multiply(p2)
	out := p.Buf()
	p.Coefficients = prod.Copy(buf).Coefficients
	return out
}

// Exp raises p to the power of n. To effiently allocate the buf it should have
// capacity of 3*(len(tc.p)*tc.pow - tc.pow + 1). A negative n gives an empty
// polynomial and n == 0 gives the constant 1. An empty polynomial to any
// positive power is empty.
func (p Poly) Exp(n int, buf []float64) Poly {
	if n < 0 {
		if cap(buf) == 0 {
			return Poly{Empty{}}
		}
		return Poly{Slice(buf[:0])}
	} else if n == 0 {
		if cap(buf) == 0 {
			return Poly{D0(1)}
		}
		return Poly{Buf(1, buf)}
	} else if n == 1 || p.Len() == 0 {
		return p.Copy(buf)
	} else if n == 2 {
		return p.Multiply(p).Copy(buf)
	}

	// https://en.wikipedia.org/wiki/Exponentiation_by_squaring
	//
	// Because of the repeated multiplication, to use the buffers efficiently,
	// a swap buffer is needed. So a total of 3 polynomials of length ln are
	// needed: sum, cur and swap.
	ln := p.Len()*n - n + 1
	s, buf := BufSplit(buf, ln)
	s = append(s, 1)
	sum := Poly{Slice(s)}

	c, buf := BufSplit(buf, ln)
	cur := p.Copy(c[:p.Len()])

	buf = BufLen(buf, ln)

	for {
		if n&1 == 1 {
			buf = sum.MultSwap(cur, buf)
		}
		n >>= 1
		if n == 0 {
			return sum
		}
		buf = cur.MultSwap(cur, buf)
	}
}

// D returns the derivative of p.
func (p Poly) D() Poly {
	return Poly{Derivative{p}}
}

// Df computes the value of p'(x).
func (p Poly) Df(x float64) float64 {
	return Poly{Derivative{p}}.F(x)
}

// Integral of the given polynomial with the constant set to c.
func (p Poly) Integral(c float64) Poly {
	return Poly{Integral{p, c}}
}

// Integral of the given polynomial with the constant set so that the value of
// Pt1(x) == y.
func (p Poly) IntegralAt(x, y float64) Poly {
	i := Integral{p, 0}
	i.C = y - Poly{i}.F(x)
	return Poly{i}
}
