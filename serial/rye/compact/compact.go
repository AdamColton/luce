package compact

import "github.com/adamcolton/luce/serial/rye"

const cmpctNil = 129

// Serializer wraps a rye.Serializer and extends it with Compact methods.
type Serializer struct {
	*rye.Serializer
}

// NewSerializer creates an instance of Serializer
func NewSerializer(size int) Serializer {
	return Serializer{&rye.Serializer{
		Size: size,
	}}
}

// MakeSerializer creates an instance of Serializer with the underlying []byte
// allocated.
func MakeSerializer(size int) Serializer {
	return Serializer{(&rye.Serializer{
		Size: size,
	}).Make()}
}

// Deserializer wraps a rye.Deserializer and extends it with Compact methods.
type Deserializer struct {
	*rye.Deserializer
}

// NewDeserializer creates an instance of Deserializer
func NewDeserializer(data []byte) Deserializer {
	return Deserializer{rye.NewDeserializer(data)}
}

// CompactUint64 writes x to the Serializer in the CompactUint64 format.
func (s Serializer) CompactUint64(x uint64) {
	if x < 129 {
		s.Byte(byte(x))
		return
	}
	idx := s.Idx
	s.Byte(cmpctNil)
	s.Data[idx] += s.Uint(0, x)
}

// CompactUint64 reads a Uint64 in compact form.
func (d Deserializer) CompactUint64() uint64 {
	b := d.Byte()
	if b < 129 {
		return uint64(b)
	}
	return d.Uint(b - cmpctNil)
}

// CompactInt64 writes an Int64 in compact form.
func (s Serializer) CompactInt64(x int64) {
	s.CompactUint64(Int64SignLSB(x))
}

// CompactInt64 reads an Int64 in compact form.
func (d Deserializer) CompactInt64() int64 {
	u := d.CompactUint64()
	sign := u & 1
	x := int64(u >> 1)
	if sign == 1 {
		x = -x
	}
	return x
}

// CompactSlice writes a byte slice in compact form.
func (s Serializer) CompactSlice(data []byte) {
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
func (d Deserializer) CompactSlice() []byte {
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
func (s Serializer) CompactString(str string) {
	s.CompactSlice([]byte(str))
}

// CompactString reads the string as a compact byte slice.
func (d Deserializer) CompactString() string {
	return string(d.CompactSlice())
}

// CompactSub returns a Sub-Deserializer where the underlying slice is from
// CompactSlice. The index of the parent is placed at the end of the data
// allocated to the Sub-Deserializer.
func (d Deserializer) CompactSub() Deserializer {
	return Deserializer{
		Deserializer: &rye.Deserializer{
			Data: d.CompactSlice(),
		},
	}
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

// SizeString returns the size of the string in compact form.
func SizeString(s string) uint64 {
	return Size([]byte(s))
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
func SizeUint64(x uint64) uint64 {
	if x < cmpctNil {
		return 1
	}
	return 1 + SizeUint(x)
}

// SizeCompactInt64 of x in compact form
func SizeInt64(x int64) uint64 {
	return SizeUint64(Int64SignLSB(x))
}

// Int64SignLSB converts an int64 to a uint64 by placing the sign in the least
// significant bit (as opposted to 2's compliment). For compact encoding, this
// makes increases the number of leading zeros that can be dropped for compact
// encoding.
func Int64SignLSB(x int64) uint64 {
	var sign uint64
	if x < 0 {
		sign = 1
		x = -x
	}
	return uint64(x<<1) | sign
}
