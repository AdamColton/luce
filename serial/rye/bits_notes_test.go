package rye_test

import (
	"testing"

	"github.com/adamcolton/luce/serial/rye"
	"github.com/stretchr/testify/assert"
)

func TestBitsWriteOrs(t *testing.T) {
	b := rye.NewBits(8)
	b.Write(1)
	b.Reset().Write(0)
	assert.Equal(t, byte(1), b.Data[0])
}

func TestBitsSubBitsIdx(t *testing.T) {
	sub := rye.NewBits(0)
	sub.WriteUint(5, 3)
	b := rye.NewBits(0)
	b.WriteSubBits(sub, 4)
	assert.Equal(t, 0, sub.Idx)
	assert.Equal(t, 7, b.Ln)

	got := b.Reset().ReadSubBits(4)
	assert.Equal(t, 3, got.Ln)
	assert.Equal(t, 3, got.Idx)
	assert.Equal(t, uint64(5), got.Reset().ReadUint(3))
}
