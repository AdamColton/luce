// Package parsers provides parsers from strings to the builtin types.
package parsers

import (
	"strconv"
)

// == projects.Code.luce.reflector ==
// [ ] parsers for all base types
//  Only String, Float64, Int, Int64 and Bool exist so far.

// String sets out to in. It never fails.
func String(out *string, in string) (err error) {
	*out = in
	return nil
}

// Float64 parses s into f. If s is not a valid float64 an error is returned and
// f is set to what strconv.ParseFloat returns, which is 0 for invalid syntax.
func Float64(f *float64, s string) (err error) {
	*f, err = strconv.ParseFloat(s, 64)
	return
}

// Int parses s as a decimal int into i. If s is not valid an error is returned
// and i is set to what strconv.Atoi returns, which is 0 for invalid syntax.
func Int(i *int, s string) (err error) {
	*i, err = strconv.Atoi(s)
	return
}

// Int64 parses s into i. The base comes from the prefix of s, as with
// strconv.ParseInt with a base of 0, so "0x10" is 16. If s is not valid an error
// is returned and i is set to what strconv.ParseInt returns.
func Int64(i *int64, s string) (err error) {
	*i, err = strconv.ParseInt(s, 0, 64)
	return
}

// Bool sets i to true if s is "y" or "Y" and to false for anything else. It
// never fails.
func Bool(i *bool, s string) (err error) {
	*i = s == "y" || s == "Y"
	return
}
