package lerr_test

import (
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/stretchr/testify/assert"
)

func TestNotEqual(t *testing.T) {
	a := 3
	b := 5
	ne := lerr.NewNotEqual(a == b, a, b)
	assert.Equal(t, "Expected 3 got 5", ne.Error())

	b = 3
	ne = lerr.NewNotEqual(a == b, a, b)
	assert.NoError(t, ne)
}
