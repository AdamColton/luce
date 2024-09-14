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

func TestSeekLock(t *testing.T) {
	seeks := []struct {
		offset int64
		whence int
	}{
		{1, io.SeekStart},
		{1, io.SeekCurrent},
		{-1, io.SeekEnd},
	}
	for _, s := range seeks {
		for _, start := range []bool{false, true} {
			// Locked, Seek leaves WriteAtIdx alone.
			b := lbuf.String("hello")
			b.LockWriteAtIdx = true
			b.WriteAtIdx = start
			b.Seek(s.offset, s.whence)
			assert.Equal(t, start, b.WriteAtIdx)

			// Unlocked, Seek sets it.
			b = lbuf.String("hello")
			b.WriteAtIdx = start
			b.Seek(s.offset, s.whence)
			assert.True(t, b.WriteAtIdx)
		}
	}
}

func TestSeekLimits(t *testing.T) {
	b := lbuf.String("hello")
	idx, err := b.Seek(-9, io.SeekCurrent)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), idx)

	idx, err = b.Seek(9, io.SeekStart)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), idx)

	idx, err = b.Seek(2, 99)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), idx)
}

func TestWriteDoesNotMoveIdx(t *testing.T) {
	b := lbuf.String("hello")
	b.Seek(1, io.SeekStart)
	n, err := b.Write([]byte("XY"))
	assert.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.Equal(t, 1, b.Idx)
	assert.Equal(t, "hXYlo", b.String())

	b.WriteAtIdx = false
	b.Write([]byte("Z"))
	assert.Equal(t, "hXYloZ", b.String())
}
