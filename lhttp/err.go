package lhttp

// Err is an error that carries an HTTP status, so ErrStatus can report it. It
// fulfills error and StatusErr.
type Err struct {
	// Code is the HTTP status.
	Code int
	// String is the error message.
	String string
}

// Error returns the message.
func (err Err) Error() string {
	return err.String
}

// Status returns the HTTP status.
func (err Err) Status() int {
	return err.Code
}
