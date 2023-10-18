package lfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExts(t *testing.T) {
	tt := map[string]struct {
		hidden   bool
		exts     []string
		expected string
	}{
		"js": {
			hidden:   false,
			exts:     []string{"js"},
			expected: "^([^\\/]|(\\/[^\\.]))*\\.((js))$",
		},
		"js-withHidden": {
			hidden:   true,
			exts:     []string{"js"},
			expected: "^.*\\.((js))$",
		},
		"js&txt": {
			hidden:   false,
			exts:     []string{"js", "txt"},
			expected: "^([^\\/]|(\\/[^\\.]))*\\.((js)|(txt))$",
		},
		"js&txt-withHidden": {
			hidden:   true,
			exts:     []string{"js", "txt"},
			expected: "^.*\\.((js)|(txt))$",
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.expected, Exts(tc.hidden, tc.exts...))
		})
	}
}
