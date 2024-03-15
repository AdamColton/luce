package lerr

import "fmt"

// ErrNotEqual is used to indicate two values that were expected to be equal.
// were not.
type ErrNotEqual struct {
	Expected, Actual interface{}
}

// NotEqual creates an instance of ErrNotEqual.
func NotEqual(expected, actual interface{}) ErrNotEqual {
	return ErrNotEqual{
		Expected: expected,
		Actual:   actual,
	}
}

// Error fulfills the error interface.
func (e ErrNotEqual) Error() string {
	return fmt.Sprintf("Expected %v got %v", e.Expected, e.Actual)
}

// NewNotEqual will return nil if areEqual is true and will create an instance
// of ErrNotEqual if it is false.
func NewNotEqual(areEqual bool, expected, actual interface{}) error {
	if areEqual {
		return nil
	}
	return NotEqual(expected, actual)
}
