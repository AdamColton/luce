package lmap_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/stretchr/testify/assert"
)

func TestDeleteManyNil(t *testing.T) {
	var w lmap.Wrapper[int, string]
	assert.NotPanics(t, func() { w.DeleteMany([]int{1, 2}) })
}
