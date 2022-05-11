package lerr

// Panic if err is not nil. If err is in the exception list, it will return
// true, but will not panic.
func Panic(err error, except ...error) bool {
	if Except(err, except...) {
		return true
	}
	if err != nil {
		panic(err)
	}
	return false
}

var LogTo func(err error)

// Log returns true if err is not nil, even if err is in the exception list. Log
// will pass the err to LogTo if it is not nil and not in the exception list.
func Log(err error, except ...error) bool {
	if Except(err, except...) {
		return false
	}
	isErr := err != nil
	if isErr && LogTo != nil {
		LogTo(err)
	}
	return isErr
}

// Except returns true if err is equal to one of the exceptions.
func Except(err error, except ...error) bool {
	for _, ex := range except {
		if err == ex {
			return true
		}
	}
	return false
}

// Any returns the first error in errs that is not nil.
func Any(errs ...error) error {
	for _, e := range errs {
		if e != nil {
			return e
		}
	}
	return nil
}
