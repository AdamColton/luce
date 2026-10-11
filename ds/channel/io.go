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

// Reader uses a []byte channel to fulfill io.Reader. Create one with a pointer,
// as it keeps what is left of the last message between calls to Read.
type Reader struct {
	Ch   <-chan []byte
	left []byte
}

// Read fulfills io.Reader. It blocks until a message is received if nothing is
// left of the last one, and a message too long for p is read over several
// calls. It returns io.EOF once the channel is closed and everything has been
// read.
func (r *Reader) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	for len(r.left) == 0 {
		msg, ok := <-r.Ch
		if !ok {
			return 0, io.EOF
		}
		r.left = msg
	}
	n = copy(p, r.left)
	r.left = r.left[n:]
	return n, nil
}

// [ ] channel.Reader.WriteTo
