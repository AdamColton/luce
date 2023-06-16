package rye

import (
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/math/ints"
)

// Bits supports reading and writing individual bits.
type Bits struct {
	Data slice.Slice[byte]
	Idx  int
	Ln   int
}

func NewBits(ln int) *Bits {
	bln := ints.DivUp(ln, 8)
	return &Bits{
		Data: make(slice.Slice[byte], bln),
		Ln:   ln,
	}
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

// ShallowCopy shares the underlying Data, but not the Ln or Idx values.
func (b *Bits) ShallowCopy() *Bits {
	return &Bits{
		Data: b.Data,
		Idx:  b.Idx,
		Ln:   b.Ln,
	}
}

// Reset the Idx to 0. Syntactic sugar.
func (b *Bits) Reset() *Bits {
	b.Idx = 0
	return b
}

// Write a single bit to Bits.
func (b *Bits) Write(bit byte) *Bits {
	idx := b.Idx / 8
	ln := idx + 1
	b.Data = b.Data.CheckCapacity(ln, 0)
	if dLn := len(b.Data); dLn < ln {
		b.Data = b.Data[:ln]
		for idx := dLn; idx < ln; idx++ {
			b.Data[idx] = 0
		}
	}
	b.Data[b.Idx/8] |= ((bit & 1) << (b.Idx % 8))
	b.Idx++
	if b.Idx > b.Ln {
		b.Ln = b.Idx
	}
	return b
}

// WriteBits takes all the bits in "from" starting at from.Idx to from.Ln and
// writes them to b. The value of from.Idx is reset to 0 after this.
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

// ReadUint of n bits.
func (b *Bits) ReadUint(n byte) uint64 {
	var u uint64
	for i := byte(0); i < n; i++ {
		u |= uint64(b.Read()) << i
	}
	return u
}

// WriteUint u as n bits.
func (b *Bits) WriteUint(u uint64, n byte) {
	for i := byte(0); i < n; i++ {
		b.Write(byte(u & 1))
		u >>= 1
	}
}

// WriteSubBits writes contents of sub with the length encoded with bitLn bits.
// This can be used to combine multiple *Bits all encoded with a common bitLn.
func (b *Bits) WriteSubBits(sub *Bits, bitLn byte) {
	b.WriteUint(uint64(sub.Ln), bitLn)
	b.WriteBits(sub.Reset())
}

// ReadSubBits reads a *Bits using bitLn to decode the length of the Bits.
func (b *Bits) ReadSubBits(bitLn byte) *Bits {
	bln := uint(b.ReadUint(bitLn))
	bi := &Bits{}
	for j := uint(0); j < bln; j++ {
		bi.Write(b.Read())
	}
	return bi
}
