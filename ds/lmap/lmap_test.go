package lmap_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/lmap/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	m := lmap.Empty[int, string](10)
	assert.Equal(t, 0, m.Len())
	m = lmap.New[int, string](nil)
	assert.Equal(t, 0, m.Len())
}

func TestEmptySafe(t *testing.T) {
	m := lmap.EmptySafe[int, string](10)
	assert.Equal(t, 0, m.Len())
	m = lmap.NewSafe[int, string](nil)
	assert.Equal(t, 0, m.Len())
}

func TestMap(t *testing.T) {
	testsuite.TestMap(lmap.New[int, string], t)
}

func TestSafe(t *testing.T) {
	testsuite.TestMap(lmap.NewSafe[int, string], t)
}
