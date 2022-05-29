package huffman

import "github.com/adamcolton/luce/serial/rye"

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
		sum.WriteBits(b)
	}

	uln := uint64(len(bs))

	s := &rye.Serializer{
		Size: int(rye.SizeCompactUint64(uln) + 1 + rye.Size(sum.Data)),
	}
	s.Make()
	s.CompactUint64(uln)
	s.Byte(bitLn)
	s.CompactSlice(sum.Data)
	return s.Data
}

func DecodeBits(data []byte) []*Bits {
	d := rye.NewDeserializer(data)
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
