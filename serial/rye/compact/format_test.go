package compact_test

import (
	"testing"

	"github.com/adamcolton/luce/serial/rye/compact"
	"github.com/stretchr/testify/assert"
)

func encode(f func(s compact.Serializer)) []byte {
	s := compact.NewSerializer(16)
	s.Make()
	f(s)
	return s.Data[:s.Idx]
}

// The format is described in the comments of CompactUint64 and CompactSlice.
func TestCompactFormat(t *testing.T) {
	assert.Equal(t, []byte{128}, encode(func(s compact.Serializer) { s.CompactUint64(128) }))
	assert.Equal(t, []byte{130, 129}, encode(func(s compact.Serializer) { s.CompactUint64(129) }))
	assert.Equal(t, []byte{131, 232, 3}, encode(func(s compact.Serializer) { s.CompactUint64(1000) }))

	assert.Equal(t, []byte{129}, encode(func(s compact.Serializer) { s.CompactSlice(nil) }))
	assert.Equal(t, []byte{5}, encode(func(s compact.Serializer) { s.CompactSlice([]byte{5}) }))
	assert.Equal(t, []byte{130, 200}, encode(func(s compact.Serializer) { s.CompactSlice([]byte{200}) }))
	assert.Equal(t, []byte{131, 1, 2}, encode(func(s compact.Serializer) { s.CompactSlice([]byte{1, 2}) }))
}
