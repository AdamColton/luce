package rye

const cmpctNil = 129

// CompactUint64 writes x to the Serializer in the CompactUint64 format.
func (s *Serializer) CompactUint64(x uint64) {
	if x < 129 {
		s.Byte(byte(x))
		return
	}
	idx := s.Idx
	s.Byte(cmpctNil)
	s.Data[idx] += s.Uint(0, x)
}

// CompactUint64 reads a Uint64 in compact form.
func (d *Deserializer) CompactUint64() uint64 {
	b := d.Byte()
	if b < 129 {
		return uint64(b)
	}
	return d.Uint(b - cmpctNil)
}

func int64ToUint64(x int64) uint64 {
	var sign uint64
	if x < 0 {
		sign = 1
		x = -x
	}
	return uint64(x<<1) | sign
}

// CompactInt64 writes an Int64 in compact form.
func (s *Serializer) CompactInt64(x int64) {
	s.CompactUint64(int64ToUint64(x))
}

// CompactInt64 reads an Int64 in compact form.
func (d *Deserializer) CompactInt64() int64 {
	u := d.CompactUint64()
	sign := u & 1
	x := int64(u >> 1)
	if sign == 1 {
		x = -x
	}
	return x
}

// CompactSlice writes a byte slice in compact form.
func (s *Serializer) CompactSlice(data []byte) {
	ln := len(data)
	if ln == 0 {
		s.Byte(cmpctNil)
		return
	}
	if ln == 1 && data[0] < cmpctNil {
		s.Byte(data[0])
		return
	}
	if ln > 121 {
		idx := s.Idx
		s.Byte(250)
		s.Data[idx] += s.Uint(0, uint64(ln))
	} else {
		s.Byte(byte(ln + cmpctNil))
	}
	s.Slice(data)
}

// CompactSlice reads a byte slice in compact form.
func (d *Deserializer) CompactSlice() []byte {
	b := d.Byte()
	if b <= cmpctNil {
		if b == cmpctNil {
			return nil
		}
		return []byte{b}
	}
	if b < 251 {
		return d.Slice(int(b - cmpctNil))
	}
	return d.Slice(int(d.Uint(b - 250)))
}

// CompactString writes the string as a compact byte slice.
func (s *Serializer) CompactString(str string) {
	s.CompactSlice([]byte(str))
}

// CompactString reads the string as a compact byte slice.
func (d *Deserializer) CompactString() string {
	return string(d.CompactSlice())
}

// Size of the data in compact form
func Size(data []byte) uint64 {
	ln := len(data)
	if ln == 0 || (ln == 1 && data[0] < cmpctNil) {
		return 1
	}
	uln := uint64(ln)
	if ln < 122 {
		return 1 + uln
	}
	return 1 + SizeUint(uln) + uln
}

// SizeUint is the number of bytes needed to encode x ignoring leading zero
// bytes.
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

// SizeCompactUint64 of x in compact form
func SizeCompactUint64(x uint64) uint64 {
	if x < cmpctNil {
		return 1
	}
	return 1 + SizeUint(x)
}

// SizeCompactInt64 of x in compact form
func SizeCompactInt64(x int64) uint64 {
	return SizeCompactUint64(int64ToUint64(x))
}
