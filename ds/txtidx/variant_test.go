package txtidx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
