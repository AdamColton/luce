package funcs

import (
	"math"
)

// Unary applies a function of one variable, F, to Of. DF is the derivative of
// F, and Name is what String calls it. Unary applies the chain rule, so any
// function of one variable with a known derivative can be used in a CDM.
type Unary struct {
	Name  string
	Of    CDM
	F, DF func(float64) float64
}

// M implements CDM.
func (u Unary) M(x []float64) float64 {
	return u.F(u.Of.M(x))
}

// IdxDM implements CDM.
func (u Unary) IdxDM(x []float64, idx int) float64 {
	// The chain rule: F(Of(x))' = F'(Of(x)) * Of'(x).
	d := u.Of.IdxDM(x, idx)
	if d == 0 {
		return 0
	}
	return u.DF(u.Of.M(x)) * d
}

// String returns the expression as Name(Of).
func (u Unary) String() string {
	return u.Name + "(" + str(u.Of) + ")"
}

// Exp returns e^of.
func Exp(of CDM) Unary {
	return Unary{Name: "exp", Of: of, F: math.Exp, DF: math.Exp}
}

// Sqrt returns the square root of of.
func Sqrt(of CDM) Unary {
	return Unary{Name: "sqrt", Of: of, F: math.Sqrt, DF: dSqrt}
}

func dSqrt(x float64) float64 {
	return 0.5 / math.Sqrt(x)
}

// Sin returns sin(of).
func Sin(of CDM) Unary {
	return Unary{Name: "sin", Of: of, F: math.Sin, DF: math.Cos}
}

// Cos returns cos(of).
func Cos(of CDM) Unary {
	return Unary{Name: "cos", Of: of, F: math.Cos, DF: dCos}
}

func dCos(x float64) float64 {
	return -math.Sin(x)
}

// Tanh returns tanh(of), which runs smoothly from -1 to 1.
func Tanh(of CDM) Unary {
	return Unary{Name: "tanh", Of: of, F: math.Tanh, DF: dTanh}
}

func dTanh(x float64) float64 {
	t := math.Tanh(x)
	return 1 - t*t
}

// Sigmoid returns the logistic function 1/(1+e^-of), which runs smoothly from
// 0 to 1. lo + (hi-lo)*Sigmoid(x) keeps a variable between lo and hi.
func Sigmoid(of CDM) Unary {
	return Unary{Name: "sigmoid", Of: of, F: sigmoid, DF: dSigmoid}
}

func sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

func dSigmoid(x float64) float64 {
	s := sigmoid(x)
	return s * (1 - s)
}

// Abs returns |of|. It isn't differentiable where of is 0, and its derivative
// there is taken as 0.
func Abs(of CDM) Unary {
	return Unary{Name: "abs", Of: of, F: math.Abs, DF: sign}
}

func sign(x float64) float64 {
	switch {
	case x > 0:
		return 1
	case x < 0:
		return -1
	}
	return 0
}
