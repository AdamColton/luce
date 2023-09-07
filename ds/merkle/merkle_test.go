package merkle_test

import (
	"crypto/rand"
	"crypto/sha256"
	"testing"

	"github.com/adamcolton/luce/ds/merkle"
	"github.com/stretchr/testify/assert"
)

func TestLeaf(t *testing.T) {
	var maxLeafSize uint32 = 200
	var ln uint32 = 5010
	data := make([]byte, ln)
	rand.Read(data)

	b := merkle.NewBuilder(maxLeafSize, sha256.New)
	m := b.Build(nil)
	assert.Nil(t, m)
	m = b.Build(data)
	assert.Equal(t, data, m.Data())

	expected := ln / maxLeafSize
	if ln%maxLeafSize != 0 {
		expected++
	}
	assert.Equal(t, int(expected), m.Leaves())
}
