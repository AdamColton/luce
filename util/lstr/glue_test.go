package lstr_test

import (
	"testing"

	"github.com/adamcolton/luce/util/lstr"
	"github.com/stretchr/testify/assert"
)

func TestGlueLengths(t *testing.T) {
	assert.Equal(t, "", lstr.Glue())
	assert.Equal(t, "one", lstr.Glue("one"))
	assert.Equal(t, "onetwothree", lstr.Glue("one", "two", "three"))
}
