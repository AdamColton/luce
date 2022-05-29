package huffman

import (
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
	data := letters
	ht := MapNew(data)
	l := NewLookup(ht)
	var bs []*Bits
	for _, r := range l.All() {
		bs = append(bs, l.Get(r))
	}

	enc := EncodeBits(bs)
	dec := DecodeBits(enc)
	assert.Equal(t, len(bs), len(dec))
	for i, b := range bs {
		assert.Equal(t, b.Ln, dec[i].Ln)
		assert.Equal(t, b.Data, dec[i].Data)
	}
}
