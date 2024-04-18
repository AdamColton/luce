package lmap_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/stretchr/testify/assert"
)

func TestPopNil(t *testing.T) {
	var w lmap.Wrapper[int, string]
	v, found := w.Pop(1)
	assert.Equal(t, "", v)
	assert.False(t, found)
}
