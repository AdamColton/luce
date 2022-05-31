package rye

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

func (b *Bits) ReadUint(n byte) uint64 {
	var u uint64
	for i := byte(0); i < n; i++ {
		u |= uint64(b.Read()) << i
	}
	return u
}

func (b *Bits) ReadUint32() uint32 {
	var u uint32
	for i := byte(0); i < 32; i++ {
		u |= uint32(b.Read()) << i
	}
	return u
}

func (b *Bits) WriteUint(u uint64, n byte) {
	for i := byte(0); i < n; i++ {
		b.Write(byte(u & 1))
		u >>= 1
	}
}

func EncodeBits(bs []*Bits) []byte {
	sum := &Bits{}
	var maxLn int
	for _, b := range bs {
		if b.Ln > maxLn {
			maxLn = b.Ln
		}

	}

	um := uint(maxLn) >> 1
	var bitLn byte = 1
	for um > 0 {
		um >>= 1
		bitLn++
	}
	for _, b := range bs {
		sum.WriteUint(uint64(b.Ln), bitLn)
	}

	for _, b := range bs {
		sum.WriteBits(b.Reset())
	}

	uln := uint64(len(bs))

	s := &Serializer{
		Size: int(SizeCompactUint64(uln) + 1 + Size(sum.Data)),
	}
	s.Make()
	s.CompactUint64(uln)
	s.Byte(bitLn)
	s.CompactSlice(sum.Data)
	return s.Data
}

func DecodeBits(data []byte) []*Bits {
	d := NewDeserializer(data)
	ln := d.CompactUint64()
	bitLn := d.Byte()
	b := &Bits{
		Data: d.CompactSlice(),
	}

	lns := make([]uint, ln)
	for i := range lns {
		lns[i] = uint(b.ReadUint(bitLn))
	}

	bs := make([]*Bits, ln)
	for i := range bs {
		bi := &Bits{}
		for j := uint(0); j < lns[i]; j++ {
			bi.Write(b.Read())
		}
		bs[i] = bi
	}
	return bs
}
