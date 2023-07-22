package lfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathLength(t *testing.T) {
	tt := map[string]string{
		"this/is/a/test.txt": "is/a/test.txt",
		"a/test.txt":         "a/test.txt",
		"/a/test.txt":        "a/test.txt",
	}

	pl := PathLength(3)
	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc, pl.Trim(n))
		})
	}
}

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
