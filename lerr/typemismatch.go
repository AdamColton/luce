package lerr

import (
	"fmt"
	"reflect"
)

// ErrTypeMismatch indicates that two types that were expected to be equal were
// not.
type ErrTypeMismatch struct {
	Expected, Actual reflect.Type
}

// TypeMismatch creates an instance of ErrTypeMismatch.
func TypeMismatch(expected, actual interface{}) ErrTypeMismatch {
	return ErrTypeMismatch{
		Expected: reflect.TypeOf(expected),
		Actual:   reflect.TypeOf(actual),
	}
}

// Error fulfills the error interface.
func (e ErrTypeMismatch) Error() string {
	return fmt.Sprintf(`Types do not match: expected "%s", got "%s"`, e.Expected.String(), e.Actual.String())
}

// NewTypeMismatch will return nil if the given types are the same and returns
// ErrTypeMismatch if they do not.
func NewTypeMismatch(expected, actual interface{}) error {
	te, ta := reflect.TypeOf(expected), reflect.TypeOf(actual)
	if te == ta {
		return nil
	}
	return ErrTypeMismatch{
		Expected: te,
		Actual:   ta,
	}
}
