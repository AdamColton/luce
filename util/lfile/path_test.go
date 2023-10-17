package lfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	tt := map[string][2]string{
		"/foo/bar.txt": {"/foo/", "bar.txt"},
		"/foo/bar/":    {"/foo/", "bar"},
		"foo.txt":      {"", "foo.txt"},
		"foo/":         {"", "foo"},
		"a":            {"", "a"},
		"/a":           {"/", "a"},
		"//a/b.txt":    {"//a/", "b.txt"},
		"/":            {"", "/"},
		"":             {"", ""},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			head, name := Name(n)
			assert.Equal(t, tc[0], head)
			assert.Equal(t, tc[1], name)
		})
	}
}
