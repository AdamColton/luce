package lerr

import "strings"

// Str provides a string error that can be set to a const making it good for
// sentinal errors.
type Str string

// Error fulfills error
func (err Str) Error() string {
	return string(err)
}

// Many allows many errors to be collected.
type Many []error

func NewMany(errs ...error) Many {
	var m Many
	for _, err := range errs {
		m = m.Add(err)
	}
	return m
}

func (m Many) First() error {
	if len(m) > 0 {
		return m[0]
	}
	return nil
}

// Error fulfills error
func (m Many) Error() string {
	out := make([]string, 0, len(m))
	for _, e := range m {
		out = append(out, e.Error())
	}
	return strings.Join(out, "\n")
}

// Add an error to the collection. If err is nil, it will not be added
func (m Many) Add(err error) Many {
	if err == nil {
		return m
	}
	return append(m, err)
}

// Cast to error type. If Many contains no errors, it will return nil.
func (m Many) Cast() error {
	if len(m) == 0 {
		return nil
	}
	return m
}
