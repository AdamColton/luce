package txtidx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoot(t *testing.T) {
	tt := map[string]string{
		"test ":      "test",
		"123test --": "123test",
		"123Test.--": "123test",
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc, root(n))
		})
	}
}
