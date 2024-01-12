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

// Must takes a value and an error. If the error is not nil, it panics. If
// the error is nil, it returns only the value.
func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

// LogTo can be set to handle errors when Log is called.
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

const (
	// ErrHandlerFunc is returned from HandlerFunc if the provided handler
	// is not func(error), chan<- error or chan error.
	ErrHandlerFunc = Str("handler argument to HandlerFunc must be func(error) or chan error")
)

// ErrHandler is a function that can handle an error.
type ErrHandler func(error)

// Handle passes err into ErrHandler if both ErrHandler and err are not nil.
// Returns a bool indicating if err was nil.
func (fn ErrHandler) Handle(err error) (isErr bool) {
	isErr = err != nil
	if fn != nil && isErr {
		fn(err)
	}
	return
}

// HandlerFunc return an ErrHandler. If the errHandler argument is an
// ErrHandler, that will be returned. If it is an error channel then that will
// be wrapped in a function and returned.
func HandlerFunc(handler any) (fn ErrHandler, err error) {
	if handler == nil {
		return
	}
	switch t := handler.(type) {
	case func(error):
		fn = t
	case chan<- error:
		fn = func(err error) {
			t <- err
		}
	case chan error:
		fn = func(err error) {
			t <- err
		}
	default:
		err = ErrHandlerFunc
	}
	return
}
