package compact_test

import (
	"math/rand"
	"testing"

	"github.com/adamcolton/luce/serial/rye"
	"github.com/adamcolton/luce/serial/rye/compact"
	"github.com/stretchr/testify/assert"
)

func TestBitsRoundTrip(t *testing.T) {
	bs := make([]*rye.Bits, 30)
	for i := range bs {
		bs[i] = &rye.Bits{}
		ln := rand.Intn(15) + 16
		for j := 0; j < ln; j++ {
			bs[i].Write(byte(rand.Intn(2)))
		}
	}

	enc := compact.EncodeBits(bs)
	dec := compact.DecodeBits(enc)
	assert.Equal(t, len(bs), len(dec))
	for i, b := range bs {
		assert.Equal(t, b.Ln, dec[i].Ln)
		assert.Equal(t, b.Data, dec[i].Data)
	}
}
