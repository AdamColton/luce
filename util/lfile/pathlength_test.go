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
