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
