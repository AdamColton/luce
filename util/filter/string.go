package filter

import (
	"regexp"
	"strings"
)

// Prefix creates a string filter that returns true when passed a string with
// the given prefix.
func Prefix(prefix string) Filter[string] {
	return func(s string) bool {
		return strings.HasPrefix(s, prefix)
	}
}

// Suffix creates a string filter that returns true when passed a string with
// the given suffix.
func Suffix(suffix string) Filter[string] {
	return func(s string) bool {
		return strings.HasSuffix(s, suffix)
	}
}

// Contains creates a string Filter that returns true when passed a string that
// contains the given substr.
func Contains(substr string) Filter[string] {
	// TODO: fix comment in previous 2 commits: filter -> Filter
	return func(s string) bool {
		return strings.Contains(s, substr)
	}
}

// Regex returns the MatchString method on regular expressions generated from
// re.
func Regex(re string) (Filter[string], error) {
	r, err := regexp.Compile(re)
	if err != nil {
		return nil, err
	}
	return r.MatchString, nil
}
