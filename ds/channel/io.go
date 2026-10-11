package channel

import (
	"io"
	"slices"
)

// chunkSize is the size of the chunks Writer.ReadFrom reads.
const chunkSize = 32 * 1024

// == projects.Code.luce.channel ==
// [ ] reduce duplication
// some of this logic overlaps with ds/bus/iobus
// is it possible to reduce the duplication?

// Writer uses a []byte channel to fulfill io.Writer
type Writer struct {
	Ch chan<- []byte
}

// Write fulfills io.Writer, sending a copy of data on the channel and blocking
// until it is received. The copy means the caller is free to reuse data as soon
// as Write returns, as the io.Writer contract requires.
func (w Writer) Write(data []byte) (n int, err error) {
	w.Ch <- slices.Clone(data)
	return len(data), nil
}

// ReadFrom fulfills io.ReaderFrom, reading from r until io.EOF and sending what
// it reads on the channel in chunks. Each chunk is a new slice that the receiver
// owns. It returns the number of bytes read and any error other than io.EOF.
func (w Writer) ReadFrom(r io.Reader) (n int64, err error) {
	for {
		buf := make([]byte, chunkSize)
		ln, err := r.Read(buf)
		if ln > 0 {
			w.Ch <- buf[:ln]
			n += int64(ln)
		}
		if err == io.EOF {
			return n, nil
		}
		if err != nil {
			return n, err
		}
	}
}

// [ ] channel.Reader.Read
// [ ] channel.Reader.WriteTo
