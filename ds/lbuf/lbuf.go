package lbuf

import (
	"io"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/math/cmpr"
	"github.com/adamcolton/luce/math/ints"
)

type Buffer struct {
	Data slice.Slice[byte]
	Idx  int
	// If WriteAtIdx is true writes will happend to the end of the Data, if it is
	// false Writes will happen at Idx.
	WriteAtIdx bool
	// If LockWriteAtIdx is false, WriteAtIdx will be set to false when doing an
	// Seek operation. If it is true, seek will have no effect on WriteAtIdx.
	LockWriteAtIdx bool
}

func String(str string) *Buffer {
	return New([]byte(str))
}

func New(buf []byte) *Buffer {
	return &Buffer{
		Data: buf,
	}
}

func (b *Buffer) Read(p []byte) (n int, err error) {
	s := b.Data[b.Idx:]
	copy(p, s)
	ln := len(s)
	if ln == 0 {
		return 0, io.EOF
	}
	ln = cmpr.Min(len(p), ln)
	b.Idx += ln
	return ln, nil
}

func (b *Buffer) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		b.Idx = int(offset)
		b.WriteAtIdx = !b.LockWriteAtIdx || b.WriteAtIdx
	case io.SeekCurrent:
		b.Idx += int(offset)
		b.WriteAtIdx = !b.LockWriteAtIdx || b.WriteAtIdx
	case io.SeekEnd:
		b.Idx = len(b.Data) + int(offset)
		b.WriteAtIdx = !b.LockWriteAtIdx || !b.WriteAtIdx
	}
	b.Idx = ints.Range(0, b.Idx, len(b.Data))
	return int64(b.Idx), nil
}

func (b *Buffer) Write(p []byte) (n int, err error) {
	lnp := len(p)
	if b.WriteAtIdx {
		idxToEnd := len(b.Data) - b.Idx
		if idxToEnd > 0 {
			copy(b.Data[b.Idx:], p)
			p = p[cmpr.Min(idxToEnd, lnp):]
		}
	}
	b.Data = append(b.Data, p...)
	return lnp, nil
}

func (b *Buffer) Len() int {
	return len(b.Data)
}

func (b *Buffer) String() string {
	return string(b.Data)
}
