package cmpr

type Tolerance float64

// DefaultTolerance adjusts how close values must be to be considered equal
var DefaultTolerance Tolerance = 1e-5

// Equal returns true if a and b are within Tolerance t of eachother.
func (t Tolerance) Equal(a, b float64) bool {
	return t.Zero(a - b)
}

// Zero returns true if x is within Tolerance t of value of 0.
func (t Tolerance) Zero(x float64) bool {
	z := Tolerance(x)
	return z < t && z > -t
}

// Equal returns true if a and b are within the DefaultTolerance of eachother.
func Equal(a, b float64) bool {
	return DefaultTolerance.Equal(a, b)
}

// Zero returns true if x is within the DefaultTolerance of 0.
func Zero(x float64) bool {
	return DefaultTolerance.Zero(x)
}

// Equalizer allows a type to define an equality test with a given Tolerance.
//
// Note that when fulfilling an interface Go will coerce a pointer type to it's
// base type, but not the other way. So if AssertEqual is on the base type
// is on the base type and a pointer to that type is passed into geomtest.Equal
// it will be cast to the base type.
type AssertEqualizer interface {
	AssertEqual(to interface{}, t Tolerance) error
}
