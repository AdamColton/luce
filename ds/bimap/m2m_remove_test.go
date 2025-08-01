package bimap_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/bimap"
	"github.com/stretchr/testify/assert"
)

func TestM2MRemoveMissing(t *testing.T) {
	m := bimap.NewM2M[string, int]()
	m.Add("a", 1)
	m.Add("c", 3)

	// Neither is in the mapping, or the pair is not present.
	m.Remove("q", 9)
	m.Remove("a", 3)
	assert.Equal(t, 2, m.LenA())
	assert.Equal(t, 2, m.LenB())

	assert.Nil(t, m.A("q"))
	assert.Nil(t, m.B(9))
	assert.Nil(t, m.RemoveA("q"))
	assert.Nil(t, m.RemoveB(9))
}
