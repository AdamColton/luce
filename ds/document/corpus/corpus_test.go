package corpus_test

import (
	"testing"

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
