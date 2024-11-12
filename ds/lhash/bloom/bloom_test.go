package bloom_test

import (
	"crypto/rand"
	"hash"
	"hash/crc32"
	"testing"

	"github.com/adamcolton/luce/ds/lhash"
	"github.com/adamcolton/luce/ds/lhash/bloom"
	"github.com/stretchr/testify/assert"
)

func TestBloom(t *testing.T) {
	idx := lhash.HashIndexer[string, uint32]{
		IndexLen: 8,
		BitLen:   8,
		Factory:  func() hash.Hash { return crc32.New(crc32.IEEETable) },
		Converter: func(t string) []byte {
			return []byte(t)
		},
	}
	f := bloom.New[string, uint32](idx)
	str := "the sun was shining on the sea"
	err := f.Add(str)
	assert.NoError(t, err)

	assert.True(t, f.Contains(str))
	assert.False(t, f.Contains(str+"FOOOO"))
}

func BenchmarkBloomAdd(b *testing.B) {
	toAdd := make([][]byte, 1000)
	for i := range toAdd {
		bs := make([]byte, 100)
		rand.Read(bs)
		toAdd[i] = bs
	}
	idx := lhash.HashIndexer[[]byte, uint32]{
		IndexLen: 8,
		BitLen:   8,
		Factory:  func() hash.Hash { return crc32.New(crc32.IEEETable) },
		Converter: func(t []byte) []byte {
			return t
		},
	}
	f := bloom.New[[]byte, uint32](idx)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		f.Add(toAdd...)
	}
}
