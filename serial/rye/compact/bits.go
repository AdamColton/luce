package compact

import "github.com/adamcolton/luce/serial/rye"

// EncodeBits takes a slice of *Bits and encodes them as a single []byte. The
// []*Bits can be recovered with DecodeBits.
func EncodeBits(bs []*rye.Bits) []byte {
	sum := &rye.Bits{}
	var maxLn int
	for _, b := range bs {
		if b.Ln > maxLn {
			maxLn = b.Ln
		}

	}

	// Find how many bits are needed to encode maxLn.
	um := uint(maxLn) >> 1
	var bitLn byte = 1
	for um > 0 {
		um >>= 1
		bitLn++
	}

	for _, b := range bs {
		sum.WriteSubBits(b, bitLn)
	}

	uln := uint64(len(bs))

	s := NewSerializer(int(SizeUint64(uln) + 1 + Size(sum.Data)))
	s.Make()
	s.CompactUint64(uln)
	s.Byte(bitLn)
	s.CompactSlice(sum.Data)
	return s.Data
}

// EncodeBits takes a slice of *Bits and encodes them as a single []byte. The
// []*Bits can be recovered with DecodeBits. The structure of the encoded data
// is a CompactUint64 for the data length. Then a byte bit length to use when
// looking up lengths. Then all the lengths are encoded, then all the Bits are
// encoded.
func DecodeBits(data []byte) []*rye.Bits {
	d := NewDeserializer(data)
	bs := make([]*rye.Bits, d.CompactUint64())
	bitLn := d.Byte()
	b := &rye.Bits{
		Data: d.CompactSlice(),
	}

	for i := range bs {
		bs[i] = b.ReadSubBits(bitLn)
	}
	return bs
}
