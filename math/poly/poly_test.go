package poly_test

import (
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/math/cmpr/cmprtest"
	"github.com/adamcolton/luce/math/poly"
	"github.com/stretchr/testify/assert"
)

func TestAtIdx(t *testing.T) {
	p := poly.New(1, 2, 3)
	cmprtest.Equal(t, 1.0, p.AtIdx(0))
	cmprtest.Equal(t, 2.0, p.AtIdx(1))
	cmprtest.Equal(t, 3.0, p.AtIdx(2))
	cmprtest.Equal(t, 0.0, p.AtIdx(3))
	cmprtest.Equal(t, 0.0, p.AtIdx(100))

	cmprtest.Equal(t, []float64{1, 2, 3}, p.Buf())

	e := poly.Poly{poly.Empty{}}
	p = poly.New()
	assert.Equal(t, p, e)
	p = poly.Poly{poly.Slice(nil)}
	cmprtest.Equal(t, p, e)
	cmprtest.Equal(t, 0.0, e.AtIdx(0))

	d0 := poly.Poly{poly.D0(5)}
	p = poly.New(5)
	assert.Equal(t, p, d0)
	p = poly.Poly{poly.Slice{5}}
	cmprtest.Equal(t, p, d0)
	cmprtest.Equal(t, 0.0, d0.AtIdx(1))

	d1 := poly.Poly{poly.D1(5)}
	p = poly.New(5, 1)
	assert.Equal(t, p, d1)
	p = poly.Poly{poly.Slice{5, 1}}
	cmprtest.Equal(t, p, d1)

	p = poly.New(1, 2, 3, 0, 0, 0)
	p2 := poly.New(1, 2, 3)
	cmprtest.Equal(t, p, p2)

	buf := make([]float64, 3)
	b := poly.Buf(3, buf)
	p = poly.New(1)
	cmprtest.Equal(t, p, poly.Poly{b})
	cmprtest.Equal(t, 1.0, buf[0])
}

func TestCopy(t *testing.T) {
	buf := make([]float64, 20)
	p := poly.New(1, 2, 3)
	cp := p.Copy(buf)
	cmprtest.Equal(t, p, cp)
	cmprtest.Equal(t, p.Buf(), buf[:3])
	cmprtest.Equal(t, 0.0, buf[4])

	cp = p.Copy(nil)
	cmprtest.Equal(t, p, cp)
}

func TestF(t *testing.T) {
	p := poly.New(5)
	cmprtest.Equal(t, 5.0, p.F(2.0))

	p = poly.New(5, 2)
	cmprtest.Equal(t, 6.0, p.F(0.5))

	p = poly.New(5, 2, 4)
	cmprtest.Equal(t, 7.0, p.F(0.5))
}

func TestAssertEqual(t *testing.T) {
	d1 := poly.Poly{poly.D1(5)}
	p := poly.Poly{poly.Slice{5, 1, 0}}

	err := d1.AssertEqual(p, 1e-10)
	assert.Nil(t, err)

	p = poly.New(1, 5)
	err = p.AssertEqual(d1, 1e-10)
	assert.Equal(t, "\t0: Expected 1 got 5\n\t1: Expected 5 got 1", err.Error())

	err = p.AssertEqual(1.0, 1e-10)
	assert.IsType(t, lerr.ErrTypeMismatch{}, err)

}

func TestDivide(t *testing.T) {
	p := poly.New(120, 154, 71, 14, 1) // (x+2)(x+3)(x+4)(x+5)
	f := 0.0

	expected := poly.New(60, 47, 12, 1)
	p, f = p.Divide(-2, p.Buf())
	cmprtest.Equal(t, expected, p)
	cmprtest.Equal(t, 0.0, f)
	assert.Equal(t, 4, p.Len())
	cmprtest.Equal(t, 0.0, p.F(-3))
	cmprtest.Equal(t, 6.0, p.F(-2))

	expected = poly.New(12, 7, 1)
	p, f = p.Divide(-5, p.Buf())
	cmprtest.Equal(t, expected, p)
	cmprtest.Equal(t, 0.0, f)
	assert.Equal(t, 3, p.Len())
	cmprtest.Equal(t, 0.0, p.F(-3))
	cmprtest.Equal(t, 2.0, p.F(-5))

	// a constant divides to nothing, and is the remainder
	p, f = poly.New(5).Divide(2, nil)
	assert.Equal(t, 0, p.Len())
	assert.Equal(t, 5.0, f)

	// an empty polynomial divides to an empty polynomial with no remainder
	p, f = poly.New().Divide(1, nil)
	assert.Equal(t, 0, p.Len())
	assert.Equal(t, 0.0, f)
}

func TestSum(t *testing.T) {
	p1 := poly.New(1, 2)
	p2 := poly.New(3, 4, 5)

	expected := poly.New(4, 6, 5)
	cmprtest.Equal(t, expected, p1.Add(p2))

	assert.Equal(t, 3, p2.Add(p1).Len())
}

func TestScale(t *testing.T) {
	got := poly.New(1, 2, 3).Scale(2)
	expected := poly.New(2, 4, 6)
	cmprtest.Equal(t, expected, got)

	got = poly.New(1, 2, 3).Scale(2)
	cmprtest.Equal(t, expected, got)
}

func TestMultiply(t *testing.T) {
	p1 := poly.New(-1, 1)
	p2 := poly.New(1, 1)

	cmprtest.Equal(t, poly.New(-1, 0, 1), p1.Multiply(p2))

	p := poly.New(1)
	p2 = poly.New(1)

	for i := 2.0; i < 6; i++ {
		x := poly.New(-i, 1)
		p = p.Multiply(x).Copy(nil)
		p2 = p2.Multiply(x)
	}
	expected := poly.New(120, -154, 71, -14, 1)
	cmprtest.Equal(t, expected, p)
	cmprtest.Equal(t, expected, p2)

	// an empty polynomial times anything is empty
	assert.Equal(t, 0, poly.New().Multiply(poly.New()).Len())
	assert.Equal(t, 0, poly.New().Multiply(p1).Len())
	assert.Equal(t, 0, p1.Multiply(poly.New()).Copy(nil).Len())
}

func TestMultSwap(t *testing.T) {
	buf, bufa, bufb := make([]float64, 10), make([]float64, 10), make([]float64, 10)
	expa := poly.New(1, 1)
	expb := poly.New(-1, 1)
	a := expa.Copy(bufa)
	b := expb.Copy(bufb)
	swap := buf

	// buf --> a
	// bufa --> swap
	swap = a.MultSwap(b, swap)
	expa = expa.Multiply(expb)
	cmprtest.Equal(t, expa, a)
	cmprtest.Equal(t, a.Buf(), buf[:3]) // a should now be in buf
	assert.Equal(t, swap, bufa[:2])     // swap will have the old value of a

	// bufa --> b
	// bufb --> swap
	swap = b.MultSwap(a, swap)
	expb = expb.Multiply(expa)
	cmprtest.Equal(t, expb, b)
	cmprtest.Equal(t, b.Buf(), bufa[:4]) // a should now be in buf
	assert.Equal(t, swap, bufb[:2])      // swap will have the old value of a

	// bufb --> a
	// buf --> swap
	swap = a.MultSwap(b, swap)
	expa = expa.Multiply(expb)
	cmprtest.Equal(t, expa, a)
	cmprtest.Equal(t, a.Buf(), bufb[:6]) // a should now be in buf
	assert.Equal(t, swap, buf[:3])       // swap will have the old value of a
}

func TestExp(t *testing.T) {
	tt := map[string]struct {
		p   poly.Poly
		pow int
	}{
		"(x2+c)^5": {
			p:   poly.New(2, -3),
			pow: 5,
		},
		"(x3+x2+c)^4": {
			p:   poly.New(1, 1, 1),
			pow: 4,
		},
		"(x4+x3+x2+c)^3": {
			p:   poly.New(4, 2, -3, 1),
			pow: 3,
		},
		"(x4+x3+x2+c)^1": {
			p:   poly.New(4, 2, -3, 1),
			pow: 1,
		},
		"(x4+x3+x2+c)^2": {
			p:   poly.New(4, 2, -3, 1),
			pow: 2,
		},
		"(x4+x3+x2+c)^0": {
			p:   poly.New(4, 2, -3, 1),
			pow: 0,
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			ln := tc.p.Len()*tc.pow - tc.pow + 1
			prod := poly.Poly{poly.Buf(ln, nil)}
			buf := make([]float64, ln)
			for i := 0; i < tc.pow; i++ {
				buf = prod.MultSwap(tc.p, buf)
			}
			buf = make([]float64, ln*3)
			cmprtest.Equal(t, prod, tc.p.Exp(tc.pow, buf))
			buf = make([]float64, ln*2+1)
			cmprtest.Equal(t, prod, tc.p.Exp(tc.pow, buf))
			buf = make([]float64, ln+1)
			cmprtest.Equal(t, prod, tc.p.Exp(tc.pow, buf))
			cmprtest.Equal(t, prod, tc.p.Exp(tc.pow, nil))

		})
	}

	// when no buffer is provided the returned value is equal to Poly{Empty{}}
	assert.Equal(t, poly.Poly{poly.Empty{}}, poly.New(4, 2, -3, 1).Exp(-1, nil))
	assert.Equal(t, poly.Poly{poly.D0(1)}, poly.New(4, 2, -3, 1).Exp(0, nil))

	// when a buffer is provided, it is used
	assert.Equal(t, poly.Poly{poly.Slice{}}, poly.New(4, 2, -3, 1).Exp(-1, []float64{1, 2, 3}))
	assert.Equal(t, poly.Poly{poly.Slice{1}}, poly.New(4, 2, -3, 1).Exp(0, []float64{5, 2, 3}))

	// an empty polynomial to a positive power is empty
	assert.Equal(t, poly.Poly{poly.D0(1)}, poly.New().Exp(0, nil))
	for _, n := range []int{1, 2, 3, 10} {
		assert.Equal(t, 0, poly.New().Exp(n, nil).Len())
	}
}

func TestD(t *testing.T) {
	cmprtest.Equal(t, poly.New(1, 8), poly.New(3, 1, 4).D())
	cmprtest.Equal(t, poly.New(1, 8, 3), poly.New(3, 1, 4, 1).D())

	p := poly.New(3, 1, 4, 1)
	d := p.D()
	cmprtest.Equal(t, poly.New(1, 8, 3), d)

	dc := poly.Poly{poly.Derivative{p}}

	for x := -10.0; x < 10.0; x += 0.1 {
		df := d.F(x)
		assert.Equal(t, df, p.Df(x))
		assert.Equal(t, df, dc.F(x))
	}

	// the derivative of a constant or an empty polynomial is empty
	assert.Equal(t, 0, poly.New(5).D().Len())
	assert.Equal(t, 0, poly.New().D().Len())
	assert.Equal(t, 0.0, poly.New().Df(2))
}
