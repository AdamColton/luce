package lstr_test

import (
	"testing"

	"github.com/adamcolton/luce/util/lstr"
	"github.com/stretchr/testify/assert"
)

func TestBufJoinSeperators(t *testing.T) {
	tt := map[string]struct {
		sep      lstr.Seperator
		elems    []string
		expected string
	}{
		"single-byte":      {",", []string{"a,", ",b", "c"}, "a,b,c"},
		"multi-byte":       {", ", []string{"a, ", ", b", "c"}, "a, b, c"},
		"multi-byte-long":  {"--", []string{"a--", "--b", "c"}, "a--b--c"},
		"empty-seperator":  {"", []string{"a", "b", "c"}, "abc"},
		"only-one-has-sep": {"--", []string{"a--", "b", "--c"}, "a--b--c"},
		"skip-empty":       {"-", []string{"a", "", "b"}, "a-b"},
	}
	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			got := tc.sep.BufJoin(tc.elems, nil)
			assert.Equal(t, tc.expected, got)
			assert.Equal(t, len(got), tc.sep.JoinLen(tc.elems))
		})
	}
}
