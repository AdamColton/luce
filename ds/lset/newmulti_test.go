package lset_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lset"
	"github.com/stretchr/testify/assert"
)

func TestNewMulti(t *testing.T) {
	a, b := lset.New(1, 2), lset.New(2, 3)
	m := lset.NewMulti(a, b)
	assert.Equal(t, lset.Multi[int]{a, b}, m)
}
