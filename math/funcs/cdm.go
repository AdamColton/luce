package funcs

import (
	"math"
	"strconv"
	"strings"

	"github.com/adamcolton/luce/ds/slice"
)

// CDM is a composable differentiable multi-func: a function of several
// variables that also gives its exact partial derivatives. CDMs are built up
// from smaller ones, such as X, Const, Sum and Product, and each works out its
// derivative from theirs.
type CDM interface {
	// M is the value of the function at x.
	M(x []float64) float64

	// IdxDM is the partial derivative of M with respect to x[idx], at x.
	IdxDM(x []float64, idx int) float64
}

// X is the variable x[i], as a CDM.
type X int

// M implements CDM.
func (i X) M(x []float64) float64 {
	return x[i]
}

// IdxDM implements CDM.
func (i X) IdxDM(x []float64, idx int) float64 {
	if i == X(idx) {
		return 1
	}
	return 0
}

// String returns the variable as x[i].
func (i X) String() string {
	return "x[" + strconv.Itoa(int(i)) + "]"
}

// Const is a constant, as a CDM.
type Const float64

// M implements CDM.
func (c Const) M([]float64) float64 {
	return float64(c)
}

// IdxDM implements CDM.
func (c Const) IdxDM(x []float64, idx int) float64 {
	return 0
}

// String returns the constant.
func (c Const) String() string {
	return formatFloat(float64(c))
}

// CoExp is C*Base^E.
type CoExp struct {
	Base CDM
	C, E float64
}

// M implements CDM.
func (ce CoExp) M(x []float64) float64 {
	return ce.C * math.Pow(ce.Base.M(x), ce.E)
}

// IdxDM implements CDM.
func (ce CoExp) IdxDM(x []float64, idx int) float64 {
	db := ce.Base.IdxDM(x, idx)
	if db == 0 {
		return 0
	}
	b := ce.Base.M(x)
	return ce.C * ce.E * db * math.Pow(b, ce.E-1)
}

// String returns the expression, leaving out a C or E of 1 and writing a C of
// -1 as a minus sign.
func (ce CoExp) String() string {
	s := wrap(ce.Base)
	if ce.E != 1 {
		s += "^" + formatFloat(ce.E)
	}
	switch ce.C {
	case 1:
	case -1:
		s = "-" + s
	default:
		s = formatFloat(ce.C) + "*" + s
	}
	return s
}

// Sum adds its terms.
type Sum []CDM

// M implements CDM.
func (s Sum) M(x []float64) (sum float64) {
	for _, cdm := range s {
		sum += cdm.M(x)
	}
	return
}

// IdxDM implements CDM.
func (s Sum) IdxDM(x []float64, idx int) (sum float64) {
	for _, cdm := range s {
		sum += cdm.IdxDM(x, idx)
	}
	return
}

// String returns the terms joined by " + ", or by " - " before a term that
// starts with a minus sign.
func (s Sum) String() string {
	terms := make([]string, len(s))
	for i, c := range s {
		terms[i] = str(c)
	}
	return joinTerms(terms)
}

// Product multiplies its factors.
type Product []CDM

// M implements CDM.
func (p Product) M(x []float64) float64 {
	prod := 1.0
	for _, cdm := range p {
		prod *= cdm.M(x)
	}
	return prod
}

// IdxDM implements CDM.
func (p Product) IdxDM(x []float64, idx int) float64 {
	// Often only one factor depends on x[idx], for example k*x[i]. Finding the
	// factors that do first means the others' derivatives are never needed
	// twice, and when none do, no factor is evaluated at all.
	nonzero, at := 0, 0
	var d float64
	for i, cdm := range p {
		if di := cdm.IdxDM(x, idx); di != 0 {
			nonzero++
			at, d = i, di
		}
	}
	switch nonzero {
	case 0:
		return 0
	case 1:
		for i, cdm := range p {
			if i != at {
				d *= cdm.M(x)
			}
		}
		return d
	}
	// The product rule for any number of factors. After each factor, prod is
	// the product of the factors so far and d its derivative:
	// (prod*v)' = prod'*v + prod*v'.
	prod := 1.0
	d = 0
	for _, cdm := range p {
		v := cdm.M(x)
		d = d*v + prod*cdm.IdxDM(x, idx)
		prod *= v
	}
	return d
}

// String returns the factors joined by "*".
func (p Product) String() string {
	return join(p, "*", true)
}

// Div is Div[0]/Div[1].
type Div [2]CDM

// M implements CDM.
func (d Div) M(x []float64) float64 {
	return d[0].M(x) / d[1].M(x)
}

// IdxDM implements CDM.
func (d Div) IdxDM(x []float64, idx int) float64 {
	dd0 := d[0].IdxDM(x, idx)
	dd1 := d[1].IdxDM(x, idx)
	if dd1 == 0 {
		if dd0 == 0 {
			return 0
		}
		return dd0 / d[1].M(x)
	}
	// The quotient rule.
	d0 := d[0].M(x)
	d1 := d[1].M(x)
	return (d1*dd0 - d0*dd1) / (d1 * d1)
}

// String returns the expression.
func (d Div) String() string {
	return wrap(d[0]) + "/" + wrap(d[1])
}

// Ln is the natural log of Of.
type Ln struct {
	Of CDM
}

// M implements CDM.
func (l Ln) M(x []float64) float64 {
	return math.Log(l.Of.M(x))
}

// IdxDM implements CDM.
func (l Ln) IdxDM(x []float64, idx int) float64 {
	dm := l.Of.IdxDM(x, idx)
	if dm == 0 {
		return 0
	}
	return dm / l.Of.M(x)
}

// String returns the expression as ln(Of).
func (l Ln) String() string {
	return "ln(" + str(l.Of) + ")"
}

// Pow is Base^Exp, where both may depend on x. Where only one does, CoExp
// (for x^k) or Exp (for e^x) is simpler.
type Pow struct {
	Base, Exp CDM
}

// M implements CDM.
func (p Pow) M(x []float64) float64 {
	return math.Pow(p.Base.M(x), p.Exp.M(x))
}

// IdxDM implements CDM.
func (p Pow) IdxDM(x []float64, idx int) float64 {
	// (u^v)' = u^v * (v'*ln(u) + v*u'/u). Each term is only computed when its
	// derivative is nonzero: ln(u) is NaN for a negative u, which would
	// otherwise make the derivative of x^2 NaN for every negative x.
	du := p.Base.IdxDM(x, idx)
	dv := p.Exp.IdxDM(x, idx)
	if du == 0 && dv == 0 {
		return 0
	}
	u := p.Base.M(x)
	v := p.Exp.M(x)
	var d float64
	if dv != 0 {
		d += dv * math.Log(u)
	}
	if du != 0 {
		d += v * du / u
	}
	return math.Pow(u, v) * d
}

// String returns the expression.
func (p Pow) String() string {
	return wrap(p.Base) + "^" + wrap(p.Exp)
}

// Linear is Coef[0]*x[0] + Coef[1]*x[1] + ... + C. It is the same as a Sum of
// CoExps, but each partial derivative is just a coefficient, so it costs
// nothing. Variables past the end of Coef have a coefficient of 0.
type Linear struct {
	Coef []float64
	C    float64
}

// M implements CDM.
func (l Linear) M(x []float64) float64 {
	sum := l.C
	for i, c := range l.Coef {
		sum += c * x[i]
	}
	return sum
}

// IdxDM implements CDM.
func (l Linear) IdxDM(x []float64, idx int) float64 {
	if idx < len(l.Coef) {
		return l.Coef[idx]
	}
	return 0
}

// String returns the expression, leaving out terms with a coefficient of 0.
func (l Linear) String() string {
	var terms []string
	for i, c := range l.Coef {
		if c != 0 {
			terms = append(terms, CoExp{Base: X(i), C: c, E: 1}.String())
		}
	}
	if l.C != 0 || len(terms) == 0 {
		terms = append(terms, formatFloat(l.C))
	}
	return joinTerms(terms)
}

// Add returns the sum of cs.
func Add(cs ...CDM) CDM {
	return Sum(cs)
}

// Sub returns a - b.
func Sub(a, b CDM) CDM {
	return Sum{a, Neg(b)}
}

// Neg returns -c.
func Neg(c CDM) CDM {
	if k, ok := c.(Const); ok {
		return -k
	}
	return Scale(-1, c)
}

// Scale returns k*c.
func Scale(k float64, c CDM) CDM {
	return CoExp{Base: c, C: k, E: 1}
}

// Square returns c².
func Square(c CDM) CDM {
	return CoExp{Base: c, C: 1, E: 2}
}

// RMSE is the root mean square of zeros. Its derivative is undefined where
// every value is zero, so to solve a system of equations, minimize SumSquares
// instead.
func RMSE(zeros ...CDM) CDM {
	ln := len(zeros)
	sum := make(Sum, ln)
	for i, z := range zeros {
		sum[i] = Square(z)
	}
	ms := Product{sum, Const(1.0 / float64(ln))}
	return Sqrt(ms)
}

// SumSquares is half the sum of the squares of zeros. It is zero exactly where
// every value is zero and is smooth everywhere, so minimizing it solves the
// system of equations zeros = 0.
func SumSquares(zeros ...CDM) CDM {
	sum := make(Sum, len(zeros))
	for i, z := range zeros {
		sum[i] = CoExp{Base: z, C: 0.5, E: 2}
	}
	return sum
}

// == projects.Code.luce.funcs ==
// [ ] reverse-mode automatic differentiation for CDM
//	BuildDM calls IdxDM once per variable, and each call evaluates the whole
//	expression again, so a gradient costs about n evaluations. Reverse mode
//	gets the whole gradient in about one pass. See also the forward-mode note
//	in _dev.md.
//	https://en.wikipedia.org/wiki/Automatic_differentiation#Reverse_accumulation

// BuildDM returns a DM that calls idxDM for each variable, such as a CDM's
// IdxDM.
func BuildDM(idxDM func(x []float64, idx int) float64) func(x, buf []float64) []float64 {
	return func(x, buf []float64) []float64 {
		buf = slice.NewBuffer(buf).Zeros(len(x))
		for i := range x {
			buf[i] = idxDM(x, i)
		}
		return buf
	}
}

// == projects.Code.luce.funcs ==
// [ ] solver-independent System
//	System returns a Descender, so a system of equations can only be solved
//	by descent. Make the system its own type, so Levenberg-Marquardt or
//	Newton's method can solve it too.

// System of equations where each CDM is equal to zero when the system is
// solved. It returns a Descender, starting at x, that minimizes SumSquares of
// the system with its exact gradient.
func System(x []float64, zeros []CDM) *Descender {
	ss := SumSquares(zeros...)
	return NewDescender(Multi{
		Ln: len(x),
		M:  ss.M,
		DM: BuildDM(ss.IdxDM),
	}, x)
}

// str returns c's String if it has one.
func str(c CDM) string {
	if s, ok := c.(interface{ String() string }); ok {
		return s.String()
	}
	return "?"
}

// wrap returns c's String, in parentheses unless c is a single term (a
// variable, a constant or a function call), so it can be used as an operand.
func wrap(c CDM) string {
	switch c.(type) {
	case X, Const, Ln, Unary, Max, Min:
		return str(c)
	}
	return "(" + str(c) + ")"
}

// join returns the Strings of cs joined by sep, wrapping each if wrapped is
// set.
func join(cs []CDM, sep string, wrapped bool) string {
	parts := make([]string, len(cs))
	for i, c := range cs {
		if wrapped {
			parts[i] = wrap(c)
		} else {
			parts[i] = str(c)
		}
	}
	return strings.Join(parts, sep)
}

// joinTerms joins the terms of a sum with " + ", or with " - " before a term
// that starts with a minus sign.
func joinTerms(terms []string) string {
	var b strings.Builder
	for i, t := range terms {
		if i > 0 {
			if rest, ok := strings.CutPrefix(t, "-"); ok {
				b.WriteString(" - ")
				t = rest
			} else {
				b.WriteString(" + ")
			}
		}
		b.WriteString(t)
	}
	return b.String()
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'g', -1, 64)
}
