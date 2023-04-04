package lstr_test

import (
	"testing"

	"github.com/adamcolton/luce/util/lstr"
	"github.com/stretchr/testify/assert"
)

func TestScanner(t *testing.T) {
	str := "Test"
	s := lstr.NewScanner(str)
	c := 0
	s.Iter().Each(func(idx int, r rune, done *bool) {
		c++
		assert.Equal(t, rune(str[idx]), r)
	})
	assert.Equal(t, c, len(str))

	s.Reset()
	assert.False(t, s.Peek(lstr.IsLower))
	assert.Equal(t, 0, s.I)
	assert.True(t, s.Match(lstr.IsUpper))
	assert.Equal(t, 1, s.I)
	assert.False(t, s.Many(lstr.IsUpper))
	assert.True(t, s.Many(lstr.IsLower))
	assert.Equal(t, 4, s.I)
	assert.True(t, s.Done())
}
