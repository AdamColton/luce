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

func TestUnpackChunks(t *testing.T) {
	strs := func(bs [][]byte) []string {
		out := make([]string, len(bs))
		for i, b := range bs {
			out[i] = string(b)
		}
		return out
	}

	p := prefix.New[uint16]()
	var stream []byte
	for _, msg := range []string{"one", "", "three"} {
		for _, part := range p.Pack([]byte(msg)) {
			stream = append(stream, part...)
		}
	}
	expected := []string{"one", "", "three"}

	// Every message in a single chunk
	assert.Equal(t, expected, strs(p.Unpack(stream)))

	// One byte at a time
	var got [][]byte
	for i := range stream {
		got = append(got, p.Unpack(stream[i:i+1])...)
	}
	assert.Equal(t, expected, strs(got))
}
