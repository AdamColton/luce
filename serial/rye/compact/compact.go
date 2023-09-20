package compact

import "github.com/adamcolton/luce/serial/rye"

// cmpctNil is the byte that means nil, or empty, for a slice and it is the base
// of the length prefixes: for a slice of up to 121 bytes the prefix is
// cmpctNil+len, and for a Uint64 it is cmpctNil+the number of bytes.
const cmpctNil = 129

// Serializer wraps a rye.Serializer and extends it with Compact methods.
type Serializer struct {
	*rye.Serializer
}

// NewSerializer creates an instance of Serializer that will write size bytes.
// It does not allocate the data, see MakeSerializer.
func NewSerializer(size int) Serializer {
	return Serializer{&rye.Serializer{
		Size: size,
	}}
}

// Deserializer wraps a rye.Deserializer and extends it with Compact methods.
type Deserializer struct {
	*rye.Deserializer
}

// NewDeserializer creates an instance of Deserializer
func NewDeserializer(data []byte) Deserializer {
	return Deserializer{rye.NewDeserializer(data)}
}

// CompactUint64 writes x to the Serializer in the CompactUint64 format. A value
// up to and including 128 is a single byte. Any other value is a byte with the
// number of bytes needed for it plus 129, followed by those bytes, little-endian
// and without leading zeros.
func (s Serializer) CompactUint64(x uint64) {
	if x < 129 {
		s.Byte(byte(x))
		return
	}
	idx := s.Idx
	s.Byte(cmpctNil)
	s.Data[idx] += s.Uint(0, x)
}

// CompactUint64 reads a Uint64 in compact form. See Serializer.CompactUint64.
func (d Deserializer) CompactUint64() uint64 {
	b := d.Byte()
	if b < 129 {
		return uint64(b)
	}
	return d.Uint(b - cmpctNil)
}

// CompactInt64 writes an Int64 in compact form, as a CompactUint64 of the value
// with the sign in the least significant bit (see Int64SignLSB). math.MinInt64
// cannot be encoded.
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

// CompactSlice writes a byte slice in compact form. A nil or empty slice is the
// single byte 129, and a slice of one byte less than 129 is that byte. A slice
// of up to 121 bytes is its length plus 129 followed by the bytes. A longer
// slice is 250 plus the number of bytes needed for its length, then the length,
// little-endian, then the bytes.
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

// CompactSlice reads a byte slice in compact form. The result shares memory
// with the data. The byte 129 is read as nil.
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

// Size of the data in compact form, as written by CompactSlice.
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

// SizeUint64 returns the size of x in CompactUint64 form.
func SizeUint64(x uint64) uint64 {
	if x < cmpctNil {
		return 1
	}
	return 1 + SizeUint(x)
}

// SizeInt64 returns the size of x in CompactInt64 form.
func SizeInt64(x int64) uint64 {
	return SizeUint64(Int64SignLSB(x))
}

// Int64SignLSB converts an int64 to a uint64 by placing the sign in the least
// significant bit (as opposed to two's complement). For compact encoding, this
// increases the number of leading zeros that can be dropped. math.MinInt64
// cannot be converted, because its absolute value overflows.
func Int64SignLSB(x int64) uint64 {
	var sign uint64
	if x < 0 {
		sign = 1
		x = -x
	}
	return uint64(x<<1) | sign
}
