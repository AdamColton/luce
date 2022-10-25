package lstr_test

import (
	"testing"

	"github.com/adamcolton/luce/util/lstr"
	"github.com/stretchr/testify/assert"
)

func TestScanner(t *testing.T) {
	str := "Test"
	idx := 0
	s := lstr.NewScanner(str)
	for ; !s.Done(); s.Next() {
		assert.Equal(t, rune(str[idx]), s.Rune)
		idx++
	}

	s.Reset()
	assert.False(t, s.Peek(lstr.IsLower))
	assert.Equal(t, 0, s.Idx)
	assert.True(t, s.Match(lstr.IsUpper))
	assert.Equal(t, 1, s.Idx)
	assert.True(t, s.Many(lstr.IsLower))
	assert.Equal(t, 4, s.Idx)
	assert.True(t, s.Done())
}
