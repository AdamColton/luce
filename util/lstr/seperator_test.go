package lstr_test

import (
	"testing"

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
			assert.Equal(t, n, s.BufJoin(tc, nil))
		})
	}

	assert.Equal(t, 0, s.JoinLen(nil))
}
