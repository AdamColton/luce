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

// TypeChecker returns a func that will check the type of actual against the
// Expected type. If actual is not of type Expected an ErrTypeMismatch is
// returned.
func TypeChecker[Expected any]() func(actual any) error {
	te := reflect.TypeOf([0]Expected{}).Elem()
	return func(actual any) error {
		ta, ok := actual.(reflect.Type)
		if !ok {
			ta = reflect.TypeOf(actual)
		}
		if te == ta {
			return nil
		}
		return ErrTypeMismatch{
			Expected: te,
			Actual:   ta,
		}
	}
}
