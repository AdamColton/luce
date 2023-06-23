package corpus_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/document"
	"github.com/adamcolton/luce/ds/document/corpus"
	"github.com/stretchr/testify/assert"
)

func TestRootID(t *testing.T) {
	c := corpus.New()
	words := []string{"this", "is", "a", "test"}
	for i, w := range words {
		assert.Equal(t, corpus.RootID(i), c.WordToID(w))
	}

	for i, w := range words {
		assert.Equal(t, w, c.IDToWord(corpus.RootID(i)))
		assert.Equal(t, corpus.RootID(i), c.WordToID(w))
	}

	assert.Equal(t, "", c.IDToWord(123))

}

func TestVariantID(t *testing.T) {
	c := corpus.New()
	vars := []document.Variant{
		{1, 32},
		{0, 32},
		{1, 13},
		{0, 13},
	}
	for i, v := range vars {
		assert.Equal(t, corpus.VariantID(i), c.VariantToID(v))
	}

	for i, v := range vars {
		id := corpus.VariantID(i)
		assert.Equal(t, v, c.IDToVariant(id))
		assert.Equal(t, id, c.VariantToID(v))
	}

	assert.Equal(t, "", c.IDToWord(123))
}
