package prefix_test

import (
	"testing"

	"github.com/adamcolton/luce/util/packeter/prefix"
	"github.com/stretchr/testify/assert"
)

func TestPrefix(t *testing.T) {
	p := prefix.New[uint32]()
	data := []byte("this is a test")
	packed := p.Pack(data)
	if assert.Len(t, packed, 2) {
		assert.Len(t, packed[0], 4)
		assert.Equal(t, packed[1], data)
	}

	unpacked := p.Unpack(packed[0])
	assert.Len(t, unpacked, 0)
	unpacked = p.Unpack(packed[1])
	if assert.Len(t, unpacked, 1) {
		assert.Equal(t, data, unpacked[0])
	}

	// creating a new copy resets p.Packer.bytes
	// exercising the full init logic on Unpacker
	p = prefix.New[uint32]()

	split := len(data) / 2
	unpacked = p.Unpack(packed[0])
	assert.Len(t, unpacked, 0)
	unpacked = p.Unpack(packed[1][:split])
	assert.Len(t, unpacked, 0)
	unpacked = p.Unpack(append(packed[1][split:], packed[0][:2]...))
	if assert.Len(t, unpacked, 1) {
		assert.Equal(t, data, unpacked[0])
	}
	unpacked = p.Unpack(append(packed[0][2:], packed[1]...))
	if assert.Len(t, unpacked, 1) {
		assert.Equal(t, data, unpacked[0])
	}
}
