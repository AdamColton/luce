package rye

// CompactUint64 writes x to the Serializer in the Compact Uint64 format.
func (s *Serializer) CompactUint64(x uint64) {
	if x < 129 {
		s.Byte(byte(x))
		return
	}
	idx := s.Idx
	s.Byte(129)
	s.Data[idx] += s.Uint(0, x)
}

func (d *Deserializer) CompactUint64() uint64 {
	b := d.Byte()
	if b < 129 {
		return uint64(b)
	}
	return d.Uint(b - 129)
}

func int64ToUint64(x int64) uint64 {
	var sign uint64
	if x < 0 {
		sign = 1
		x = -x
	}
	return uint64(x<<1) | sign
}

func (s *Serializer) CompactInt64(x int64) {
	s.CompactUint64(int64ToUint64(x))
}

func (d *Deserializer) CompactInt64() int64 {
	u := d.CompactUint64()
	sign := u & 1
	x := int64(u >> 1)
	if sign == 1 {
		x = -x
	}
	return x
}

func (s *Serializer) PrefixSlice(data []byte) {
	ln := len(data)
	if ln == 0 {
		s.Byte(129)
		return
	}
	if ln == 1 && data[0] < 129 {
		s.Byte(data[0])
		return
	}
	if ln > 121 {
		idx := s.Idx
		s.Byte(250)
		s.Data[idx] += s.Uint(0, uint64(ln))
	} else {
		s.Byte(byte(ln + 129))
	}
	s.Slice(data)
}

func (d *Deserializer) PrefixSlice() []byte {
	b := d.Byte()
	if b == 129 {
		return nil
	}
	if b < 129 {
		return []byte{b}
	}
	if b < 251 {
		return d.Slice(int(b - 129))
	}
	return d.Slice(int(d.Uint(b - 250)))
}

func (s *Serializer) PrefixString(str string) {
	s.PrefixSlice([]byte(str))
}

func (d *Deserializer) PrefixString() string {
	return string(d.PrefixSlice())
}

func Size(data []byte) uint64 {
	ln := len(data)
	if ln == 0 || (ln == 1 && data[0] < 129) {
		return 1
	}
	uln := uint64(ln)
	if ln < 122 {
		return 1 + uln
	}
	return 1 + SizeUint(uln) + uln
}

func SizeUint(x uint64) uint64 {
	if x == 0 {
		return 1
	}
	var out uint64
	for ; x > 0; x >>= 8 {
		out++
	}
	return out
}

func SizeCompactUint64(x uint64) uint64 {
	if x < 129 {
		return 1
	}
	return 1 + SizeUint(x)
}

func SizeCompactInt64(x int64) uint64 {
	return SizeCompactUint64(int64ToUint64(x))
}
