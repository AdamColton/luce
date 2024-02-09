package lstr

import "strings"

// Len is a strongly typed version of the builtin len for strings. It is useful
// when a func value is needed.
func Len(s string) int {
	return len(s)
}

// NewRemover creates a strings.Replacer that removes every occurrence of the
// given strings.
func NewRemover(rm ...string) *strings.Replacer {
	s := make([]string, len(rm)*2)
	for i, r := range rm {
		s[i*2] = r
		s[i*2+1] = ""
	}
	return strings.NewReplacer(s...)
}
