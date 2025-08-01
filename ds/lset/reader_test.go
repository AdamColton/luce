package lset_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lset"
	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	var r lset.Reader[int] = lset.New(1, 2, 3)
	assert.True(t, r.Contains(2))
	assert.Equal(t, 3, r.Len())
	assert.Equal(t, 3, r.Copy().Len())
	assert.ElementsMatch(t, []int{1, 2, 3}, r.Slice(nil))
}
