package filter

import "strings"

// Prefix creates a String filter that returns true when passed a string with
// the given prefix.
func Prefix(prefix string) String {
	return func(s string) bool {
		return strings.HasPrefix(s, prefix)
	}
}

// Suffix creates a String filter that returns true when passed a string with
// the given suffix.
func Suffix(suffix string) String {
	return func(s string) bool {
		return strings.HasSuffix(s, suffix)
	}
}

// Contains creates a String filter that returns true when passed a string that
// contains the given substr.
func Contains(substr string) String {
	return func(s string) bool {
		return strings.Contains(s, substr)
	}
}
