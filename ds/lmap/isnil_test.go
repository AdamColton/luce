package lmap_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/stretchr/testify/assert"
)

func TestIsNil(t *testing.T) {
	var w lmap.Wrapper[int, string]
	assert.True(t, w.IsNil())
	assert.False(t, lmap.New(map[int]string{}).IsNil())
}
