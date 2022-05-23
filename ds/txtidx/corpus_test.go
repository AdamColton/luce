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

func TestVariant(t *testing.T) {
	tt := []string{
		"Test",
		"TEST",
		"tesTIng. ",
	}

	for _, n := range tt {
		t.Run(n, func(t *testing.T) {
			rt := root(n)
			v := findVariant(rt, n)
			assert.Equal(t, n, v.apply(rt))
		})
	}
}
