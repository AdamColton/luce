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
