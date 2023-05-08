package document_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/document"
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
			assert.Equal(t, tc, document.Root(n))
		})
	}
}

func TestVariant(t *testing.T) {
	tt := []string{
		"Test",
		"TEST",
		"tesTIng. ",
	}

	for _, word := range tt {
		t.Run(word, func(t *testing.T) {
			rt, v := document.RootVariant(word)
			assert.Equal(t, word, string(v.Apply(rt, nil)))
		})
	}
}
