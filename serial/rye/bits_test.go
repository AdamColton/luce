package rye

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitUint(t *testing.T) {
	b := &Bits{}
	b.WriteUint(12, 4)
	b.Reset()
	u := b.ReadUint(4)
	assert.Equal(t, uint64(12), u)
}

func TestBitsRoundTrip(t *testing.T) {
	bs := make([]*Bits, 30)
	for i := range bs {
		bs[i] = &Bits{}
		ln := rand.Intn(15) + 16
		for j := 0; j < ln; j++ {
			bs[i].Write(byte(rand.Intn(2)))
		}
	}

	enc := EncodeBits(bs)
	dec := DecodeBits(enc)
	assert.Equal(t, len(bs), len(dec))
	for i, b := range bs {
		assert.Equal(t, b.Ln, dec[i].Ln)
		assert.Equal(t, b.Data, dec[i].Data)
	}
}

func TestBitsCopy(t *testing.T) {
	bs := &Bits{}
	bs.WriteUint(5, 5)
	bs.WriteUint(123, 5)
	assert.Equal(t, []byte{101, 3}, bs.Data)
	assert.Equal(t, 10, bs.Idx)
	assert.Equal(t, 10, bs.Ln)

	cp := bs.Copy()
	assert.Equal(t, []byte{101, 3}, cp.Data)
	assert.Equal(t, 10, cp.Idx)
	assert.Equal(t, 10, cp.Ln)

	bs.Data[0] = 1
	cp.Data[0] = 2
	assert.Equal(t, []byte{1, 3}, bs.Data)
	assert.Equal(t, []byte{2, 3}, cp.Data)
}

func TestShallowCopy(t *testing.T) {
	b := &Bits{
		Data: []byte{3, 1},
		Ln:   12,
		Idx:  2,
	}
	cp := b.ShallowCopy()
	assert.Equal(t, b, cp)

	b.Data[0] = 10
	assert.Equal(t, b, cp)

	b.Idx = 5
	b.Ln = 15
	assert.NotEqual(t, b, cp)
}
