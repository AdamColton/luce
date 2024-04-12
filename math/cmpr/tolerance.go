package cmpr

// Tolerance represents a difference below which two floats can be considered
// equal.
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

// Unique takes a sorted list and returns all the unique values where two
// values within Tolerance of eachother are considered equal.
func (t Tolerance) Unique(s []float64) []float64 {
	if len(s) < 2 {
		return s
	}
	cur := 0
	next := 1
	for next < len(s) {
		if t.Equal(s[cur], s[next]) {
			next++
		} else {
			cur++
			s[cur] = s[next]
			next++
		}
	}
	return s[:cur+1]
}

// Equal returns true if a and b are within the DefaultTolerance of eachother.
func Equal(a, b float64) bool {
	return DefaultTolerance.Equal(a, b)
}

// Zero returns true if x is within the DefaultTolerance of 0.
func Zero(x float64) bool {
	return DefaultTolerance.Zero(x)
}
