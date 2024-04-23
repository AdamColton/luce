package lstr_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/util/lstr"
	"github.com/stretchr/testify/assert"
)

func TestSeperatorJoin(t *testing.T) {
	tt := map[string][]string{
		"":        {},
		"a":       {"a"},
		"b/c":     {"b", "c"},
		"d/e":     {"d", "", "e"},
		"f/g":     {"f/", "g"},
		"h/i":     {"h", "/i"},
		"j/k":     {"j/", "/k"},
		"/l/m/n/": {"/l/", "/m/", "/n/"},
		"o/p":     {"o/", "/", "/p"},
	}

	s := lstr.Seperator("/")
	for n, tc := range tt {
		t.Run("_"+n, func(t *testing.T) {
			assert.Equal(t, n, s.Join(tc...))
		})
	}

	assert.Equal(t, 0, s.JoinLen(nil))
}

func TestSeperatorIndex(t *testing.T) {
	tt := map[string]int{
		"":       -1,
		"*":      0,
		"a*b":    1,
		"a*b*":   1,
		"abcd*e": 4,
	}

	s := lstr.Seperator("*")
	for n, tc := range tt {
		t.Run("_"+n, func(t *testing.T) {
			assert.Equal(t, tc, s.Index(n))
		})
	}
}

func TestSeperatorSplit(t *testing.T) {
	tt := map[string]slice.Slice[string]{
		"":       {""},
		"/":      {"", ""},
		"a":      {"a"},
		"b/c":    {"b", "c"},
		"/d/e":   {"", "d", "e"},
		"fg/h/i": {"fg", "h", "i"},
	}

	s := lstr.Seperator("/")
	for n, tc := range tt {
		t.Run("_"+n, func(t *testing.T) {
			assert.Equal(t, tc, s.Split(n))
		})
	}

	assert.Equal(t, 0, s.JoinLen(nil))
}

func TestSeperatorStrings(t *testing.T) {
	str := "this\nis\na\ntest"
	strs := lstr.NewLine.Strings(str)
	assert.Equal(t, strs.Strings, []string{"this", "is", "a", "test"})
}

func TestJoiner(t *testing.T) {
	s := lstr.Seperator(";")
	j := s.Joiner("this", "is", "a", "test")
	assert.Equal(t, "this;is;a;test", j.String())
}
