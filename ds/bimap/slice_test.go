package bimap_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/bimap"
	"github.com/stretchr/testify/assert"
)

func TestSlice(t *testing.T) {
	s := bimap.NewBiSlice[string](nil, nil)

	a := s.Upsert("a")
	assert.Equal(t, a, s.Upsert("a"))
	assert.Equal(t, 1, s.Upsert("b"))
}
