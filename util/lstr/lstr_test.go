package lstr_test

import (
	"testing"

	"github.com/adamcolton/luce/util/lstr"
	"github.com/stretchr/testify/assert"
)

func TestLen(t *testing.T) {
	str := "testing"
	assert.Equal(t, len(str), lstr.Len(str))
}

func TestGlue(t *testing.T) {
	got := lstr.Glue("This", "is", "a", "test")
	assert.Equal(t, "Thisisatest", got)
}
