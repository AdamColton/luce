package lbuf

import (
	"io"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/math/cmpr"
	"github.com/adamcolton/luce/math/ints"
)

type Buffer struct {
	Data    slice.Slice[byte]
	ReadIdx int
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
	s := b.Data[b.ReadIdx:]
	copy(p, s)
	ln := cmpr.Min(len(p), len(s))
	if ln == 0 {
		return 0, io.EOF
	}
	b.ReadIdx += ln
	return ln, nil
}

func (b *Buffer) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		b.ReadIdx = int(offset)
	case io.SeekCurrent:
		b.ReadIdx += int(offset)
	case io.SeekEnd:
		b.ReadIdx = len(b.Data) + int(offset)
	}
	b.ReadIdx = ints.Range(0, b.ReadIdx, len(b.Data))
	return int64(b.ReadIdx), nil
}

func (b *Buffer) Write(p []byte) (n int, err error) {
	b.Data = append(b.Data, p...)
	return len(p), nil
}

func (b *Buffer) Len() int {
	return len(b.Data)
}
