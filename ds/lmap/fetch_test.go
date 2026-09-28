package lmap_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/stretchr/testify/assert"
)

func TestFetch(t *testing.T) {
	calls := 0
	gen := func() string {
		calls++
		return "generated"
	}

	m := lmap.New(map[int]string{1: "one"})
	assert.Equal(t, "one", m.Fetch(1, gen))
	assert.Equal(t, 0, calls)

	// A missing value is generated and stored.
	assert.Equal(t, "generated", m.Fetch(2, gen))
	assert.Equal(t, 1, calls)
	assert.Equal(t, "generated", m.Fetch(2, gen))
	assert.Equal(t, 1, calls)
}

func TestFetchNil(t *testing.T) {
	var w lmap.Wrapper[int, string]
	calls := 0
	gen := func() string {
		calls++
		return "generated"
	}

	// Nothing is stored, so the generator runs every time.
	assert.Equal(t, "generated", w.Fetch(1, gen))
	assert.Equal(t, "generated", w.Fetch(1, gen))
	assert.Equal(t, 2, calls)
	assert.True(t, w.IsNil())
}
