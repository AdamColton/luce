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

// Unique takes a sorted list and returns the unique values, where two values
// within Tolerance of each other are considered equal. Each value is compared
// to the last unique value kept, so a run of values that are each within
// Tolerance of the next can span more than Tolerance. The values are moved down
// in s, so s is modified and the result shares its memory.
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

// Equal returns true if a and b are within the DefaultTolerance of each other.
func Equal(a, b float64) bool {
	return DefaultTolerance.Equal(a, b)
}

// Zero returns true if x is within the DefaultTolerance of 0.
func Zero(x float64) bool {
	return DefaultTolerance.Zero(x)
}
