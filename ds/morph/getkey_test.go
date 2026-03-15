package morph_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/morph"
	"github.com/stretchr/testify/assert"
)

func TestGetKey(t *testing.T) {
	get := morph.GetKey[string, int]()
	assert.Equal(t, "a", get("a", 1))
	assert.Equal(t, "b", get("b", 2))
}
