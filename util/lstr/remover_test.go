package lstr_test

import (
	"testing"

	"github.com/adamcolton/luce/util/lstr"
	"github.com/stretchr/testify/assert"
)

func TestNewRemover(t *testing.T) {
	rm := lstr.NewRemover(",", "$", "ab")
	assert.Equal(t, "1234.5 c", rm.Replace("$1,234.5 abc"))
	assert.Equal(t, "unchanged", lstr.NewRemover().Replace("unchanged"))
}
