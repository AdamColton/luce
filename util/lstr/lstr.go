package lstr

import (
	"strings"

	"github.com/adamcolton/luce/math/ints"
)

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

// addLen returns ln+add and panics if that overflows an int.
func addLen(ln, add int) int {
	if add > ints.MaxI-ln {
		panic("lstr: Glue output length overflow")
	}
	return ln + add
}

// Glue strings together with no joining string. Equivalent to
// strings.Join(strs, ""). It panics if the length of the result would overflow
// an int.
func Glue(strs ...string) string {
	switch len(strs) {
	case 0:
		return ""
	case 1:
		return strs[0]
	}

	var ln int
	for _, s := range strs {
		ln = addLen(ln, len(s))
	}
	out := make([]byte, 0, ln)
	for _, s := range strs {
		out = append(out, s...)
	}
	return string(out)
}
