package lstr_test

import (
	"regexp"
	"testing"

	"github.com/adamcolton/luce/util/lstr"
	"github.com/stretchr/testify/assert"
)

func TestStringsRegexNoMatch(t *testing.T) {
	re := regexp.MustCompile(`z(\d)`)

	// With skipEmpty, every remaining string is tried and then it is done.
	s := lstr.NewStrings([]string{"a", "b"})
	assert.Nil(t, s.Regex(re, true))
	assert.True(t, s.Done())

	// Without it, only the current string is tried.
	s = lstr.NewStrings([]string{"a", "z1"})
	assert.Empty(t, s.Regex(re, false))
	assert.False(t, s.Done())
	assert.Equal(t, []string{"z1", "1"}, s.Regex(re, false))
}
