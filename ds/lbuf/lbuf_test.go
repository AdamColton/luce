package lbuf_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lbuf"
	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	b := lbuf.New([]byte("testing"))
	cp := make([]byte, 4)
	n, err := b.Read(cp)
	assert.NoError(t, err)
	assert.Equal(t, 4, n)
	assert.Equal(t, []byte("test"), cp)

	n, err = b.Read(cp)
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, []byte("ing"), cp[:n])
}
