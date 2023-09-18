package filter

import "strings"

// Prefix creates a String filter that returns true when passed a string with
// the given prefix.
func Prefix(prefix string) Filter[string] {
	return func(s string) bool {
		return strings.HasPrefix(s, prefix)
	}
}
