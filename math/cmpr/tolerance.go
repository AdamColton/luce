package cmpr

// Tolerance represents a difference below which two floats can be considered
// equal.
type Tolerance float64

// DefaultTolerance is the Tolerance used by Equal and Zero. It adjusts how close
// values must be to be considered equal.
var DefaultTolerance Tolerance = 1e-5

// Equal returns true if a and b are within Tolerance t of each other.
func (t Tolerance) Equal(a, b float64) bool {
	return t.Zero(a - b)
}

// Zero returns true if x is within Tolerance t of 0.
func (t Tolerance) Zero(x float64) bool {
	z := Tolerance(x)
	return z < t && z > -t
}

// Equal returns true if a and b are within the DefaultTolerance of each other.
func Equal(a, b float64) bool {
	return DefaultTolerance.Equal(a, b)
}

// Zero returns true if x is within the DefaultTolerance of 0.
func Zero(x float64) bool {
	return DefaultTolerance.Zero(x)
}
