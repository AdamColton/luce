package lerr

// Str provides a string error that can be set to a const making it good for
// sentinal errors.
type Str string

// Error fulfills error
func (err Str) Error() string {
	return string(err)
}
