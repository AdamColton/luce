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
