package compact_test

import (
	"bytes"
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

func encodeSlice(data []byte) []byte {
	s := compact.NewSerializer(int(compact.Size(data)))
	s.Make()
	s.CompactSlice(data)
	return s.Data
}

// The format is described in the comments of CompactUint64 and CompactSlice.
func TestCompactFormat(t *testing.T) {
	assert.Equal(t, []byte{128}, encode(func(s compact.Serializer) { s.CompactUint64(128) }))
	assert.Equal(t, []byte{130, 129}, encode(func(s compact.Serializer) { s.CompactUint64(129) }))
	assert.Equal(t, []byte{131, 232, 3}, encode(func(s compact.Serializer) { s.CompactUint64(1000) }))

	assert.Equal(t, []byte{254}, encodeSlice(nil))
	assert.Equal(t, []byte{200}, encodeSlice([]byte{200}))
	assert.Equal(t, []byte{255, 254}, encodeSlice([]byte{254}))
	assert.Equal(t, []byte{255, 0, 1, 2}, encodeSlice([]byte{1, 2}))

	header := func(ln int) []byte {
		enc := encodeSlice(make([]byte, ln))
		return enc[:len(enc)-ln]
	}
	assert.Equal(t, []byte{255, 246}, header(248))
	assert.Equal(t, []byte{255, 247, 0}, header(249))
	assert.Equal(t, []byte{255, 247, 255}, header(504))
	assert.Equal(t, []byte{255, 248, 0, 0}, header(505))
	assert.Equal(t, []byte{255, 248, 255, 255}, header(66040))
	assert.Equal(t, []byte{255, 249, 0, 0, 0}, header(66041))
}

// Every byte sequence decodes to exactly one slice, and encoding that slice
// gives back the bytes that were read. So no slice has two encodings.
func FuzzCompactSliceCanonical(f *testing.F) {
	for _, seed := range [][]byte{{}, {7}, {254}, {255, 255}, {255, 3, 1, 2, 3, 4, 5}, {255, 247, 9}} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) == 0 {
			return
		}
		// The longest header is 9 bytes, so padding means only the slice body
		// can be cut short. A read that reaches the end may have been clamped.
		padded := append(data, make([]byte, 9)...)
		d := compact.NewDeserializer(padded)
		got := d.CompactSlice()
		if d.Done() {
			return
		}
		read := padded[:d.Idx]
		assert.True(t, bytes.Equal(read, encodeSlice(got)), "%v", read)
	})
}
