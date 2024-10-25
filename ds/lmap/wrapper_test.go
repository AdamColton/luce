package lmap_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/stretchr/testify/assert"
)

func TestWrapped(t *testing.T) {
	m := lmap.New(map[int]string{1: "1", 2: "2", 3: "3"})
	_, ok := m.Mapper.(lmap.Map[int, string])
	assert.True(t, ok)

	w := lmap.Wrap(m)
	_, ok = w.Mapper.(lmap.Map[int, string])
	assert.True(t, ok)

	assert.Equal(t, w.Mapper, w.Wrapped())
}
