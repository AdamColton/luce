// Package prefix frames messages by prefixing each with its length.
package prefix

import (
	"unsafe"

	"golang.org/x/exp/constraints"
)

// Packer prefixes each message with its length. The length is written as a U,
// least significant byte first, so U decides how many bytes the prefix takes and
// the longest message that can be packed. The zero value is ready to use. It
// fulfills packeter.Packer.
type Packer[U constraints.Unsigned] struct {
	bytes int
}

func (p *Packer[U]) setBytes() {
	var u U
	p.bytes = int(unsafe.Sizeof(u))
}

// Pack returns the length of data as a U, followed by data. They are returned as
// two slices that are meant to be written in order. If the length does not fit
// in a U it is truncated.
func (p *Packer[U]) Pack(data []byte) [][]byte {
	if p.bytes == 0 {
		p.setBytes()
	}
	ln := U(len(data))
	lnbs := make([]byte, p.bytes)
	for i := range lnbs {
		lnbs[i] = byte(ln >> (8 * i))
	}
	return [][]byte{
		lnbs,
		data,
	}
}

// Unpacker fulfills packeter.Unpacker for messages that were prefixed by a
// Packer with the same U. It keeps an unfinished message between calls to
// Unpack, so the data can arrive in chunks of any size. It embeds a Packer, so
// it can pack as well as unpack. It is not safe for concurrent use.
type Unpacker[U constraints.Unsigned] struct {
	buf, lnbs []byte
	ln        U
	*Packer[U]
}

// New creates an Unpacker for messages with a length prefix of type U.
func New[U constraints.Unsigned]() *Unpacker[U] {
	return &Unpacker[U]{
		Packer: &Packer[U]{},
	}
}

func (p *Unpacker[U]) setBytes() {
	if p.bytes == 0 {
		p.Packer.setBytes()
	}
	p.lnbs = make([]byte, 0, p.bytes)
}

// Unpack takes the next chunk of data and returns every message that has now
// been completed, which may be none or several. A message that is not complete
// is kept for the next call.
func (p *Unpacker[U]) Unpack(data []byte) [][]byte {
	sizeSet := p.ln > 0
	if !sizeSet {
		data, sizeSet = p.setSize(data)
	}
	var out [][]byte
	if sizeSet {
		bytesNeeded := p.ln - U(len(p.buf))
		if U(len(data)) >= bytesNeeded {
			out = append(out, append(p.buf, data[:bytesNeeded]...))
			data = data[bytesNeeded:]
			p.buf = nil
			p.ln = 0

			if len(data) > 0 {
				out = append(out, p.Unpack(data)...)
			}
		} else {
			p.buf = append(p.buf, data...)
		}
	}
	return out
}

func (p *Unpacker[U]) setSize(data []byte) ([]byte, bool) {
	if cap(p.lnbs) == 0 {
		p.setBytes()
	}
	bytesNeeded := p.bytes - len(p.lnbs)
	sizeSet := len(data) >= bytesNeeded
	if sizeSet {
		p.lnbs = append(p.lnbs, data[:bytesNeeded]...)
		for i := p.bytes - 1; i >= 0; i-- {
			p.ln = (p.ln << 8) | U(p.lnbs[i])
		}
		p.lnbs = p.lnbs[:0]
		data = data[bytesNeeded:]
	} else {
		p.lnbs = append(p.lnbs, data...)
		data = nil
	}
	return data, sizeSet
}
