package lmap_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/stretchr/testify/assert"
)

func TestEnsure(t *testing.T) {
	var w lmap.Wrapper[int, string]
	w = w.Ensure()
	assert.False(t, w.IsNil())
	w.Set(1, "a")

	// A Wrapper that is not nil is returned as is.
	assert.Equal(t, "a", w.Ensure().GetVal(1))
}
