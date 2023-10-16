package ints_test

import (
	"testing"

	"github.com/adamcolton/luce/math/ints"
	"github.com/stretchr/testify/assert"
)

func TestDiv(t *testing.T) {
	assert.Equal(t, 3, ints.DivUp(5, 2))
	assert.Equal(t, 2, ints.DivDown(5, 2))
}
