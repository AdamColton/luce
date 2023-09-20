package rye_test

import (
	"testing"

	"github.com/adamcolton/luce/serial/rye"
	"github.com/stretchr/testify/assert"
)

func TestSerializerUintAdvancesIdx(t *testing.T) {
	s := &rye.Serializer{Size: 8}
	s.Make()

	assert.Equal(t, byte(2), s.Uint(0, 0x0102))
	assert.Equal(t, 2, s.Idx)
	assert.Equal(t, byte(3), s.Uint(3, 5))
	assert.Equal(t, 5, s.Idx)
	assert.Equal(t, []byte{2, 1, 5, 0, 0, 0, 0, 0}, s.Data)
}
