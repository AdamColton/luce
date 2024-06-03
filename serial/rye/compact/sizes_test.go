package compact_test

import (
	"testing"

	"github.com/adamcolton/luce/serial/rye/compact"
	"github.com/stretchr/testify/assert"
)

func TestMakeSerializer(t *testing.T) {
	s := compact.MakeSerializer(4)
	assert.Len(t, s.Data, 4)
	assert.Equal(t, 4, s.Size)

	// NewSerializer does not allocate.
	assert.Len(t, compact.NewSerializer(4).Data, 0)
}

func TestSizeString(t *testing.T) {
	assert.Equal(t, uint64(1), compact.SizeString(""))
	assert.Equal(t, uint64(4), compact.SizeString("abc"))
	long := make([]byte, 300)
	assert.Equal(t, compact.Size(long), compact.SizeString(string(long)))
}
