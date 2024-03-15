package lerr

import (
	"fmt"
)

// ErrLenMismatch represents a mis-matched length.
type ErrLenMismatch struct {
	Expected, Actual int
}

// LenMismatch returns an ErrLenMismatch.
func LenMismatch(expected, actual int) ErrLenMismatch {
	return ErrLenMismatch{
		Expected: expected,
		Actual:   actual,
	}
}

// Error fulfills the error interface.
func (e ErrLenMismatch) Error() string {
	return fmt.Sprintf("Lengths do not match: Expected %v got %v", e.Expected, e.Actual)
}

// NewLenMismatch will return a nil error if expected and actual are equal and
// an instance of LenMismatch if they are different. The min and max will be set
// to the smaller and larger of the values passed in respectivly.
func NewLenMismatch(expected, actual int) (min, max int, err error) {
	min, max = expected, actual

	if expected == actual {
		return
	}
	if actual < expected {
		min, max = actual, expected
	}
	err = LenMismatch(expected, actual)
	return
}
