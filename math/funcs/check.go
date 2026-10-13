package funcs

import (
	"math"

	"github.com/adamcolton/luce/ds/slice"
)

// DefaultCheckTol is the tolerance to pass to CheckDM and DMAgrees when there
// is no reason to choose another. It is well above the error of NumericDM and
// well below the error of a wrong derivative.
const DefaultCheckTol = 1e-6

// DMCheck is what CheckDM found for one partial derivative.
type DMCheck struct {
	// DM is the value from Multi.DM, and Numeric the value from NumericDM.
	DM, Numeric float64
	// Precise is the value from PrecisePartial, and ErrEst its error estimate.
	// Both are zero unless Numeric disagreed with DM.
	Precise, ErrEst float64
	// Pass is true when DM agrees with the numeric estimates.
	Pass bool
}

// CheckDM compares s.DM, which must be set, against numeric estimates of each
// partial derivative at x, and returns what it found for each, using buf if it
// has the capacity. It is for testing an analytic gradient.
//
// A partial derivative passes if it is within tol of NumericDM, relative to
// max(1, |value|). If not, NumericDM may be the one that is wrong, so CheckDM
// falls back to PrecisePartial, and the partial passes if it is within tol of
// that, allowing for its error estimate. It also fails if that error estimate
// is larger than tol, because the precise value can't settle the question.
// Ridders describes where the estimates are unreliable, and x should stay away
// from points where f isn't differentiable.
func (s *Multi) CheckDM(x []float64, tol float64, buf []DMCheck) []DMCheck {
	out := slice.NewBuffer(buf).Zeros(s.Ln)
	dm := s.DM(x, nil)
	for i := range out {
		c := &out[i]
		c.DM = dm[i]
		c.Numeric = s.M.PartialDerivative(x, i)
		if relDiff(c.DM, c.Numeric) <= tol {
			c.Pass = true
			continue
		}
		c.Precise, c.ErrEst = s.M.PrecisePartial(x, i)
		scale := math.Max(1, math.Max(math.Abs(c.DM), math.Abs(c.Precise)))
		allowed := tol * scale
		c.Pass = c.ErrEst <= allowed && math.Abs(c.DM-c.Precise) <= allowed+c.ErrEst
	}
	return out
}

// DMAgrees runs CheckDM and reports whether every partial derivative passed.
// If not, idx is the first that failed; otherwise it is -1.
func (s *Multi) DMAgrees(x []float64, tol float64) (ok bool, idx int) {
	for i, c := range s.CheckDM(x, tol, nil) {
		if !c.Pass {
			return false, i
		}
	}
	return true, -1
}

// relDiff is the difference between a and b, relative to the larger of them
// once that is above 1, so it is an absolute difference for small values.
func relDiff(a, b float64) float64 {
	return math.Abs(a-b) / math.Max(1, math.Max(math.Abs(a), math.Abs(b)))
}
