package prefix

import (
	"unsafe"

	"golang.org/x/exp/constraints"
)

type Packer[U constraints.Unsigned] struct {
	bytes int
}

func (p *Packer[U]) setBytes() {
	var u U
	p.bytes = int(unsafe.Sizeof(u))
}

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

type Unpacker[U constraints.Unsigned] struct {
	buf, lnbs []byte
	ln        U
	*Packer[U]
}

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
