package huffman

// Bits supports reading and writing individual bits.
type Bits struct {
	Data []byte
	Idx  int
	Ln   int
}

// Copy the Bits.
func (b *Bits) Copy() *Bits {
	out := &Bits{
		Data: make([]byte, len(b.Data)),
		Idx:  b.Idx,
		Ln:   b.Ln,
	}
	copy(out.Data, b.Data)
	return out
}

// Reset the Idx to 0. Syntactic sugar.
func (b *Bits) Reset() *Bits {
	b.Idx = 0
	return b
}

// Write a single bit to Bits.
func (b *Bits) Write(bit byte) *Bits {
	idx := b.Idx / 8
	for idx >= len(b.Data) {
		b.Data = append(b.Data, 0)
	}
	b.Data[b.Idx/8] |= ((bit & 1) << (b.Idx % 8))
	b.Idx++
	if b.Idx > b.Ln {
		b.Ln = b.Idx
	}
	return b
}

// WriteBits takes all the bits in from starting at from.Idx to from.Ln and
// writes them to b.
func (b *Bits) WriteBits(from *Bits) *Bits {
	for from.Idx < from.Ln {
		b.Write(from.Read())
	}
	from.Reset()
	return b
}

// Read a single bit.
func (b *Bits) Read() byte {
	bit := (b.Data[b.Idx/8] >> (b.Idx % 8)) & 1
	b.Idx++
	return bit
}
