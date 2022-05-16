package lhttp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
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

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, n, Join(tc...))
		})
	}
}
