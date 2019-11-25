package thresher

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestField(t *testing.T) {
	f := field{
		name: "Name",
		kind: 123,
	}
	assert.Equal(t, f, deserializeField(f.serialize()))
	assert.Equal(t, uint64(11736191855920939007), f.id())
}
