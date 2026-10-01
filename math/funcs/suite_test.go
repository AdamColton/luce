package funcs_test

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/adamcolton/luce/math/funcs"
	"github.com/stretchr/testify/assert"
)

// The solver suite is a set of minimization problems with known answers. Each
// solver runs every problem and the results are compared by steps, function
// evaluations and how close the answer is. TestSolverSuite logs the table (run
// with -v) and fails if a problem that a solver is known to solve stops
// converging.

// problem is a function to minimize, where to start and the known answer.
type problem struct {
	name  string
	f     funcs.M
	start []float64
	// min is the known minimizer, or nil when it isn't unique.
	min  []float64
	fmin float64
}

// residuals is a system of equations: each value is zero at the solution.
type residuals func(x []float64) []float64

// sumSquares returns ½‖r(x)‖², the objective for solving a system of
// equations. It is smooth, and for a linear system it is a quadratic.
func sumSquares(r residuals) funcs.M {
	return func(x []float64) float64 {
		sum := 0.0
		for _, v := range r(x) {
			sum += v * v
		}
		return sum / 2
	}
}

// converged reports whether x is close enough to the problem's answer: within
// a relative 1e-6 of min when it is known, or within 1e-12 of fmin.
func (p problem) converged(x []float64) bool {
	if p.min == nil {
		return p.f(x)-p.fmin <= 1e-12
	}
	for i, m := range p.min {
		if math.Abs(x[i]-m) > 1e-6*math.Max(1, math.Abs(m)) {
			return false
		}
	}
	return true
}

func problems() []problem {
	quadratic := func(scales ...float64) funcs.M {
		return func(x []float64) (sum float64) {
			for i, s := range scales {
				sum += s * x[i] * x[i]
			}
			return
		}
	}

	// Series circuit: a 10V source from node 0 to node 2, 100Ω from node 0 to
	// node 1 and 200Ω from node 1 to node 2. The unknowns are the node voltages
	// and the source current: v0, v1, v2, i. Node 0 is ground.
	series := residuals(func(x []float64) []float64 {
		v0, v1, v2, i := x[0], x[1], x[2], x[3]
		i01 := (v0 - v1) / 100
		i12 := (v1 - v2) / 200
		return []float64{
			-i01 + i,     // current into node 0
			i01 - i12,    // current into node 1
			i12 - i,      // current into node 2
			v2 - v0 - 10, // the source
			v0,           // ground
		}
	})

	// Divider: 5V through 100Ω to node n, then 200Ω and 300Ω in parallel to
	// ground. The unknowns are the source current and the voltage at n.
	divider := residuals(func(x []float64) []float64 {
		i, vn := x[0], x[1]
		return []float64{
			i + vn/100 - 5.0/100,
			5.0/100 - vn/100 - vn/200 - vn/300,
		}
	})

	// Wide range: 10V through 1Ω to node n, then 1MΩ to ground. The unknowns
	// are the voltage at n and the source current.
	wide := residuals(func(x []float64) []float64 {
		vn, i := x[0], x[1]
		return []float64{
			i - (10 - vn),
			(10 - vn) - vn/1e6,
		}
	})
	wideV := 10 * 1e6 / (1e6 + 1)

	// Diode: 5V through 1kΩ to node v, then a diode to ground. Nonlinear and
	// very steep once the diode conducts.
	diode := residuals(func(x []float64) []float64 {
		v := x[0]
		return []float64{(5-v)/1000 - 1e-12*(math.Exp(v/0.025)-1)}
	})

	// Line fit: y = a*x + b through points that aren't on a line, so the
	// minimum is above zero.
	px := []float64{0, 1, 2, 3}
	py := []float64{1, 2, 2, 4}
	fit := residuals(func(x []float64) []float64 {
		r := make([]float64, len(px))
		for i := range px {
			r[i] = x[0]*px[i] + x[1] - py[i]
		}
		return r
	})
	// The least-squares line through those points is y = 0.9x + 0.9.
	fitMin := []float64{0.9, 0.9}

	return []problem{
		{name: "bowl", f: quadratic(1, 1), start: []float64{3, -4}, min: []float64{0, 0}},
		{name: "bowl-cond1e3", f: quadratic(1, 1e3), start: []float64{3, -4}, min: []float64{0, 0}},
		{name: "bowl-cond1e6", f: quadratic(1, 1e6), start: []float64{3, -4}, min: []float64{0, 0}},
		{
			name:  "rosenbrock",
			f:     sumSquares(func(x []float64) []float64 { return []float64{10 * (x[1] - x[0]*x[0]), 1 - x[0]} }),
			start: []float64{-1.2, 1},
			min:   []float64{1, 1},
		},
		// The two circles of the existing tests cross, so there is more than one
		// point where the distance is zero.
		{name: "circles", f: dist, start: []float64{0.4, 0.3}},
		{name: "series", f: sumSquares(series), start: []float64{0, 0, 10, 0}, min: []float64{0, 10.0 / 3, 10, -1.0 / 30}},
		{name: "divider", f: sumSquares(divider), start: []float64{0.02, 2}, min: []float64{5.0 / 220, 5 * 120.0 / 220}},
		{name: "wide-range", f: sumSquares(wide), start: []float64{0, 0}, min: []float64{wideV, 10 - wideV}},
		{name: "diode", f: sumSquares(diode), start: []float64{0}},
		{name: "line-fit", f: sumSquares(fit), start: []float64{0, 0}, min: fitMin, fmin: sumSquares(fit)(fitMin)},
	}
}

// result is what one solver did on one problem.
type result struct {
	converged bool
	// steps is the step where it first converged, or the steps it took.
	steps int
	// evals counts the calls the solver made to the function.
	evals int
	x     []float64
	gap   float64
}

// solver runs a problem. f counts the evaluations; converged is for the
// solver to stop as soon as it reaches the answer, and doesn't count.
type solver struct {
	name string
	// solves lists the problems the solver is known to solve. The suite fails
	// if one of them stops converging.
	solves []string
	run    func(p problem, f funcs.M, converged func([]float64) bool) (x []float64, steps int)
}

func (s solver) expected(name string) bool {
	for _, n := range s.solves {
		if n == name {
			return true
		}
	}
	return false
}

const maxSteps = 20000

var solvers = []solver{
	{
		// The Descender as it is today, with the defaults from Init.
		name:   "descender",
		solves: []string{"bowl", "line-fit"},
		run: func(p problem, f funcs.M, converged func([]float64) bool) ([]float64, int) {
			d := (&funcs.Descender{
				Multi: funcs.Multi{Ln: len(p.start), M: f},
				Steps: maxSteps,
				X:     append([]float64(nil), p.start...),
			}).Init()
			for step := 1; d.Steps > 0; step++ {
				d.Step()
				if converged(d.X) {
					return d.X, step
				}
			}
			return d.X, maxSteps
		},
	},
}

func solve(s solver, p problem) result {
	evals := 0
	counted := func(x []float64) float64 {
		evals++
		return p.f(x)
	}
	x, steps := s.run(p, counted, p.converged)
	return result{
		converged: p.converged(x),
		steps:     steps,
		evals:     evals,
		x:         x,
		gap:       p.f(x) - p.fmin,
	}
}

func TestSolverSuite(t *testing.T) {
	var table strings.Builder
	fmt.Fprintf(&table, "\n%-12s %-14s %-9s %7s %10s %12s\n", "solver", "problem", "converged", "steps", "evals", "f - fmin")
	for _, s := range solvers {
		for _, p := range problems() {
			r := solve(s, p)
			fmt.Fprintf(&table, "%-12s %-14s %-9t %7d %10d %12.3g\n", s.name, p.name, r.converged, r.steps, r.evals, r.gap)
			if s.expected(p.name) {
				assert.True(t, r.converged, "%s no longer solves %s: x = %v", s.name, p.name, r.x)
			}
		}
	}
	t.Log(table.String())
}

func BenchmarkSolverSuite(b *testing.B) {
	for _, s := range solvers {
		for _, p := range problems() {
			b.Run(s.name+"/"+p.name, func(b *testing.B) {
				var r result
				for b.Loop() {
					r = solve(s, p)
				}
				b.ReportMetric(float64(r.evals), "evals")
				b.ReportMetric(float64(r.steps), "steps")
			})
		}
	}
}

// TestProblems checks that each problem's known answer is right.
func TestProblems(t *testing.T) {
	for _, p := range problems() {
		if p.min != nil {
			assert.InDelta(t, p.fmin, p.f(p.min), 1e-20, p.name)
		}
	}
}
