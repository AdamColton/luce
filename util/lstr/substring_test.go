package lstr_test

import (
	"testing"

	"github.com/adamcolton/luce/util/lstr"
	"github.com/stretchr/testify/assert"
)

func TestSubStrings(t *testing.T) {
	subs := lstr.SubStringBySplit([]int{4, 6, 7, 11})
	assert.Equal(t, lstr.SubStrings{{0, 4}, {4, 6}, {6, 7}, {7, 11}}, subs)
	assert.Equal(t, []string{"This", "Is", "A", "Test"}, subs.Slice("ThisIsATest"))
}
