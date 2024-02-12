package poly

import (
	"math"

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

// Newton's method to find one root of the polynomial. The initial guess is
// passed in as x; min sets how close to 0 is acceptible and it will return if a
// value closer than that is found; steps limits the maximum number of
// iterations that will; d is the derivative. It is not required to provide d,
// but if there is a cached instance available, it reduces repeated computation.
func (p Poly) Newton(x float64, min cmpr.Tolerance, steps int, d Coefficients) (float64, float64) {
	const (
		small cmpr.Tolerance = 1e-5
	)

	if d == nil {
		d = Derivative{p}
	}
	dp := Poly{d}

	y := p.F(x)

	bestY, bestX := math.Abs(y), x
	for i := 0; i < steps && !min.Zero(y); i++ {
		if math.IsInf(y, 0) || math.IsNaN(y) {
			x += 1e-3
			continue
		}
		d := dp.F(x)
		if small.Zero(d) {
			x += 1e-3
			y = p.F(x)
			continue
		}
		d = y / d
		d *= (200 - float64(i)) / 200
		x -= d
		y = p.F(x)
		if absy := math.Abs(y); absy < bestY {
			bestX, bestY = x, absy
		}
	}
	return bestX, bestY
}

// Halley's method to find one root of the polynomial. The initial guess is
// passed in as x; min sets how close to 0 is acceptible and it will return if a
// value closer than that is found; steps limits the maximum number of
// iterations that will; d is the derivative; d2 is the second derivative. It is
// not required to provide d or d2, but if there is a cached instance available,
// it reduces repeated computation.
func (p Poly) Halley(x float64, min cmpr.Tolerance, steps int, d, d2 Coefficients) (float64, float64) {
	const (
		small cmpr.Tolerance = 1e-5
	)

	if d == nil {
		d = Derivative{p}
	}
	if d2 == nil {
		d2 = Derivative{d}
	}
	dp, ddp := Poly{d}, Poly{d2}

	y := p.F(x)
	bestY, bestX := math.Abs(y), x

	for i := 0; i < steps && !min.Zero(y); i++ {
		dy := dp.F(x)
		d2y := ddp.F(x)
		denom := 2*dy*dy - y*d2y
		if small.Zero(denom) {
			x += 1e-3
			y = p.F(x)
			continue
		}
		d := (2 * y * dy) / denom
		d *= (200 - float64(i)) / 200
		x -= d
		y = p.F(x)
		if x == bestX {
			x += 1e-3
		} else if absy := math.Abs(y); absy < bestY {
			bestX, bestY = x, absy
		}
	}
	return bestX, bestY
}

// Quad finds the real roots of a quadratic equation. The number of roots to
// return is set by the length of the buffer. If the length is zero then the max
// number of roots will be found.
func Quad(c, b, a float64, buf []float64) []float64 {
	outLn := len(buf)
	if a == 0 {
		if b == 0 {
			return nil
		}
		return append(buf[:0], -c/b)
	}

	s := b*b - 4*a*c
	if s < 0 {
		return nil
	}
	if s == 0 {
		return append(buf[:0], -b/(2*a))
	}
	s = math.Sqrt(s)
	a *= 2
	buf = append(buf[:0], (-b+s)/(a))
	if outLn != 1 {
		buf = append(buf, (-b-s)/(a))
	}
	return buf
}

const (
	third float64 = 1.0 / 3.0
	sqrt3         = 1.732050807568877293527446341505872366942805253810380628055806
)

// Cubic finds the real roots of a cubic equation. The number of roots to return
// is set by the length of the buffer. If the length is zero then the max number
// of roots will be found.
func Cubic(d, c, b, a float64, buf []float64) []float64 {
	if a == 0 {
		return Quad(d, c, b, buf)
	}
	outLn := len(buf)
	if outLn == 0 {
		outLn = 3
	}

	//https://github.com/shril/CubicEquationSolver/blob/master/CubicEquationSolver.py
	a2 := a * a
	b2 := b * b

	f := ((3 * c / a) - (b2 / a2)) / 3

	a3 := a2 * a
	b3 := b2 * b
	g := (2*b3/a3 - 9*b*c/a2 + 27*d/a) / 27

	g2 := g * g
	f3 := f * f * f
	h := g2/4 + f3/27

	if f == 0 && g == 0 && h == 0 {
		return append(buf, -powThird(d/a))
	}

	var z0, z1, z2 float64
	if h <= 0 {
		i := math.Sqrt(g2/4 - h)
		j := math.Pow(i, third)
		k := math.Acos(-g/(2*i)) / 3
		L := -j
		M := math.Cos(k)
		N := sqrt3 * math.Sin(k)
		P := -b / (3 * a)

		z0 = 2*j*math.Cos(k) - (b / (3 * a))
		z1 = L*(M+N) + P
		z2 = L*(M-N) + P

	} else {

		srh := math.Sqrt(h)
		g = -g / 2
		a *= 3

		i := powThird(g + srh)
		j := powThird(g - srh)

		z0 = (i + j) - (b / (a))

		// one real root, the other two are a complex pair
		z1 = math.NaN()
		z2 = z1
	}

	buf = append(buf[:0], z0)
	if z0 != z1 && !math.IsNaN(z1) && outLn > 1 {
		buf = append(buf, z1)
	}
	if z0 != z2 && z1 != z2 && !math.IsNaN(z2) && outLn > 2 {
		buf = append(buf, z2)
	}
	return buf
}

func powThird(x float64) float64 {
	if x >= 0 {
		return math.Pow(x, third)
	}
	return -math.Pow(-x, third)
}

// Quartic finds the real roots of a Quartic equation. The number of roots to
// return is set by the length of the buffer. If the length is zero then the max
// number of roots will be found. Roots that are within 1e-7 of each other are
// taken to be one root that is repeated, and are returned once.
func Quartic(e, d, c, b, a float64, buf []float64) []float64 {
	// https://stackoverflow.com/a/50747781
	if a == 0 {
		return Cubic(e, d, c, b, buf)
	}
	outLn := len(buf)
	if outLn == 0 {
		outLn = 4
	}

	b /= a
	c /= a
	d /= a
	e /= a

	// Depress the quartic: with x = y - b/4 it is y^4 + p*y^2 + q*y + r.
	var out []float64
	b2 := b * b
	p := c - 0.375*b2
	b3 := b2 * b
	q := 0.125*b3 - 0.5*b*c + d
	r := e - 0.25*b*d + 0.0625*b2*c - 0.01171875*b3*b
	shift := -0.25 * b

	if q == 0.0 {
		// Biquadratic: y^4 + p*y^2 + r is a quadratic in z = y^2, and each
		// z >= 0 gives the roots y = +-sqrt(z).
		for _, z := range Quad(r, p, 1, nil) {
			sqrt_z, ok := sqrtTol(z, 1+math.Abs(p)+math.Abs(r))
			if !ok {
				continue
			}
			out = quarticAppend(out, outLn, shift+sqrt_z)
			out = quarticAppend(out, outLn, shift-sqrt_z)
		}
		return out
	}

	// Ferrari's method. m is the largest root of the resolvent cubic, which is
	// positive when q != 0.
	m := quarticM(p, 0.25*p*p-r, -0.125*q*q)
	sqrt_2m := math.Sqrt(2.0 * m)
	qs := q / sqrt_2m
	scale := 1 + math.Abs(m) + math.Abs(p) + math.Abs(qs)
	if delta, ok := sqrtTol(2.0*(-m-p+qs), scale); ok {
		out = quarticAppend(out, outLn, 0.5*(-sqrt_2m+delta)+shift)
		out = quarticAppend(out, outLn, 0.5*(-sqrt_2m-delta)+shift)
	}

	if delta, ok := sqrtTol(2.0*(-m-p-qs), scale); ok {
		out = quarticAppend(out, outLn, 0.5*(sqrt_2m+delta)+shift)
		out = quarticAppend(out, outLn, 0.5*(sqrt_2m-delta)+shift)
	}

	return out
}

// sqrtTol is the square root of v. A repeated root of the quartic makes v 0, but
// rounding leaves it a little above or below, so a v within 1e-12*scale of 0 is
// taken to be 0. ok is false if v is negative.
func sqrtTol(v, scale float64) (sqrt float64, ok bool) {
	tol := 1e-12 * scale
	if v < -tol {
		return 0, false
	}
	if v <= tol {
		return 0, true
	}
	return math.Sqrt(v), true
}

func quarticAppend(out []float64, outLn int, r float64) []float64 {
	const zero cmpr.Tolerance = 1e-7
	if len(out) == outLn {
		return out
	}
	for _, o := range out {
		if zero.Equal(r, o) {
			return out
		}
	}
	return append(out, r)
}

// quarticM is the largest real root of the cubic m^3 + b*m^2 + c*m + d.
func quarticM(b, c, d float64) float64 {
	// depress it: with m = t - b/3 it is t^3 + p*t + q
	p := c - b*b/3.0
	q := 2.0*b*b*b/27.0 - b*c/3.0 + d

	if p == 0.0 {
		return -math.Cbrt(q) - b/3.0
	}
	if q == 0.0 {
		if p < 0.0 {
			return math.Sqrt(-p) - b/3.0
		}
		return -b / 3.0
	}

	t := math.Sqrt(math.Abs(p) / 3.0)
	g := 1.5 * q / (p * t)
	if p > 0.0 {
		return -2.0*t*math.Sinh(math.Asinh(g)/3.0) - b/3.0
	}

	if 4.0*p*p*p+27.0*q*q < 0.0 {
		return 2.0*t*math.Cos(math.Acos(clamp(g))/3.0) - b/3.0
	}
	if q > 0.0 {
		return -2.0*t*math.Cosh(math.Acosh(math.Max(1, -g))/3.0) - b/3.0
	}
	return 2.0*t*math.Cosh(math.Acosh(math.Max(1, g))/3.0) - b/3.0
}

// clamp limits g to [-1, 1], for acos, which rounding can push just outside.
func clamp(g float64) float64 {
	return math.Max(-1, math.Min(1, g))
}
