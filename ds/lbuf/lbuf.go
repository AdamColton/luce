// Package lbuf provides a byte buffer over a slice that supports reading,
// writing and seeking.
package lbuf

import (
	"io"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/math/cmpr"
	"github.com/adamcolton/luce/math/ints"
)

// Buffer is a byte buffer over a slice. It fulfills io.Reader, io.Writer and
// io.Seeker. Reads and Seeks move Idx, Writes do not.
type Buffer struct {
	// Data holds all the bytes in the Buffer.
	Data slice.Slice[byte]
	// Idx is the position of the next read. If WriteAtIdx is true, it is also
	// where the next write happens.
	Idx int
	// If WriteAtIdx is true, writes overwrite Data starting at Idx and anything
	// beyond the end of Data is appended. If it is false, writes are appended to
	// the end of Data.
	WriteAtIdx bool
	// If LockWriteAtIdx is false, WriteAtIdx will be set to true when doing a
	// Seek operation. If it is true, Seek will have no effect on WriteAtIdx.
	LockWriteAtIdx bool
}

// String creates a Buffer from a string.
func String(str string) *Buffer {
	return New([]byte(str))
}

// New creates a Buffer that uses buf as its Data. The slice is not copied.
func New(buf []byte) *Buffer {
	return &Buffer{
		Data: buf,
	}
}

// Read fulfills io.Reader. It reads from Idx and advances Idx by the number of
// bytes read. It returns io.EOF if there is no data left to read.
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

// Seek fulfills io.Seeker. The resulting Idx is limited to the range of Data,
// so a position before the start becomes 0 and a position past the end becomes
// the length of Data, rather than an error. An unknown whence leaves Idx
// unchanged. Unless LockWriteAtIdx is set, Seek sets WriteAtIdx to true.
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
		b.WriteAtIdx = !b.LockWriteAtIdx || b.WriteAtIdx
	}
	b.Idx = ints.Range(0, b.Idx, len(b.Data))
	return int64(b.Idx), nil
}

// Write fulfills io.Writer. If WriteAtIdx is true, p overwrites Data starting
// at Idx and anything beyond the end of Data is appended, otherwise p is
// appended to the end of Data. Idx does not change. It never returns an error.
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

// Len returns the length of all of Data, not just the part after Idx.
func (b *Buffer) Len() int {
	return len(b.Data)
}

// String returns all of Data as a string, regardless of Idx.
func (b *Buffer) String() string {
	return string(b.Data)
}
