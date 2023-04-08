package lset_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lset"
	"github.com/stretchr/testify/assert"
)

func TestNilSet(t *testing.T) {
	var s *lset.Set[int]
	assert.Equal(t, 0, s.Len())
	assert.Nil(t, s.Slice(nil))
}
