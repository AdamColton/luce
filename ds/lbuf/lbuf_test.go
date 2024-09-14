package lbuf_test

import (
	"io"
	"testing"

	"github.com/adamcolton/luce/ds/lbuf"
	"github.com/stretchr/testify/assert"
)

func TestBuffer(t *testing.T) {
	b := lbuf.String("testing")
	assert.Equal(t, b.Len(), len(b.Data))

	cp := make([]byte, 0)
	n, err := b.Read(cp)
	assert.Equal(t, 0, n)
	assert.NoError(t, err)

	cp = make([]byte, 4)
	n, err = b.Read(cp)
	assert.NoError(t, err)
	assert.Equal(t, 4, n)
	assert.Equal(t, []byte("test"), cp)

	n, err = b.Read(cp)
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, []byte("ing"), cp[:n])

	n, err = b.Read(cp)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, 0, n)

	b.Idx = 4
	b.Seek(2, io.SeekStart)
	b.Seek(2, io.SeekCurrent)
	b.Write([]byte(" case"))
	assert.Equal(t, "test case", b.String())

	b.Seek(-2, io.SeekEnd)
	n, err = b.Read(cp)
	assert.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.Equal(t, []byte("se"), cp[:n])
}

func TestBufferWriteSeek(t *testing.T) {
	b := lbuf.String("test")
	b.Write([]byte("ing"))
	cp := make([]byte, 7)
	b.Read(cp)
	assert.Equal(t, "testing", string(cp))
}
