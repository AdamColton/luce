package rye_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/serial/rye"
	"github.com/stretchr/testify/assert"
)

func TestBitUint(t *testing.T) {
	b := rye.NewBits(2)
	b.WriteUint(12, 4)
	b.Reset()
	u := b.ReadUint(4)
	assert.Equal(t, uint64(12), u)
}

func TestBitsCopy(t *testing.T) {
	bs := &rye.Bits{}
	bs.WriteUint(5, 5)
	bs.WriteUint(123, 5)
	expected := slice.New([]byte{101, 3})
	assert.Equal(t, expected, bs.Data)
	assert.Equal(t, 10, bs.Idx)
	assert.Equal(t, 10, bs.Ln)

	cp := bs.Copy()
	assert.Equal(t, expected, cp.Data)
	assert.Equal(t, 10, cp.Idx)
	assert.Equal(t, 10, cp.Ln)

	bs.Data[0] = 1
	cp.Data[0] = 2
	expected = slice.New([]byte{1, 3})
	assert.Equal(t, expected, bs.Data)
	expected = slice.New([]byte{2, 3})
	assert.Equal(t, expected, cp.Data)
}

func TestWriteBits(t *testing.T) {
	to := &rye.Bits{}
	to.WriteUint(31, 5)
	expected := slice.New([]byte{31})
	assert.Equal(t, expected, to.Data)

	from := &rye.Bits{}
	from.WriteUint(31, 5)

	toRef := to.WriteBits(from.Reset())
	assert.Equal(t, to, toRef)

	expected = slice.New([]byte{255, 3})
	assert.Equal(t, expected, to.Data)
}

func TestSubBitsRoundTrip(t *testing.T) {
	acc := &rye.Bits{}
	for i := uint64(0); i < 100; i++ {
		sub := &rye.Bits{}
		sub.WriteUint(i, 8)
		acc.WriteSubBits(sub, 4)
	}
	acc.Reset()
	for i := uint64(0); i < 100; i++ {
		sub := acc.ReadSubBits(4)
		expected := &rye.Bits{}
		expected.WriteUint(i, 8)
		assert.Equal(t, expected.Data, sub.Data)
	}
}
