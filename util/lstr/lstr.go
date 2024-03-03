package lstr

import (
	"strings"

	"github.com/adamcolton/luce/math/ints"
)

// Len creates a strongly typed version of builtin len for strings.
func Len(s string) int {
	return len(s)
}

// NewRemover creates a strings.Replacer that replaces all the strings
// given with "".
func NewRemover(rm ...string) *strings.Replacer {
	s := make([]string, len(rm)*2)
	for i, r := range rm {
		s[i*2] = r
		s[i*2+1] = ""
	}
	return strings.NewReplacer(s...)
}

// Glue strings together with no joining string. Equivalent to
// strings.Join(strs, "")
func Glue(strs ...string) string {
	switch len(strs) {
	case 0:
		return ""
	case 1:
		return strs[0]
	}

	var ln int
	for _, s := range strs {
		if len(s) > ints.MaxI-ln {
			panic("lstr: Glue output length overflow")
		}
		ln += len(s)
	}
	out := make([]byte, 0, ln)
	for _, s := range strs {
		out = append(out, s...)
	}
	return string(out)
}
