package lmap_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/stretchr/testify/assert"
)

func TestLen(t *testing.T) {
	m := lmap.Map[int, string]{1: "1", 2: "2", 3: "3"}
	assert.Equal(t, len(m), m.Len())
}
