package compact

import (
	"github.com/adamcolton/luce/serial/rye"
)

// uintPrefix is the base of the CompactUint64 prefix: a value of n bytes is
// prefixed with uintPrefix+n.
const uintPrefix = 129

// The first byte of a CompactSlice is either the single byte of the slice
// (below sliceNil), sliceNil, or sliceExt. After sliceExt, a byte below
// sliceMeta is a short length; sliceMeta and above select a tier.
const (
	sliceNil  = 254
	sliceExt  = 255
	sliceMeta = 247
	shortMin  = 2
)

// tierStart is the shortest length in each tier. Tier t follows sliceMeta+t
// and stores ln-tierStart[t] in t+1 bytes, so each tier starts one past the
// end of the one before and every length has exactly one encoding.
var tierStart = [...]uint64{
	249,
	505,
	66041,
	16843257,
	4311810553,
	1103823438329,
	282578800148985,
}

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

// CompactUint64 writes x to the Serializer in the CompactUint64 format. A value
// up to and including 128 is a single byte. Any other value is a byte with the
// number of bytes needed for it plus 129, followed by those bytes, little-endian
// and without leading zeros.
func (s Serializer) CompactUint64(x uint64) {
	if x < uintPrefix {
		s.Byte(byte(x))
		return
	}
	idx := s.Idx
	s.Byte(uintPrefix)
	s.Data[idx] += s.Uint(0, x)
}

// CompactUint64 reads a Uint64 in compact form. See Serializer.CompactUint64.
func (d Deserializer) CompactUint64() uint64 {
	b := d.Byte()
	if b < uintPrefix {
		return uint64(b)
	}
	return d.Uint(b - uintPrefix)
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

// CompactSlice writes a byte slice in compact form. Every slice has exactly one
// encoding:
//
//	[]                  254
//	[b], b < 254        b
//	[b], b is 254, 255  255 b
//	len 2 to 248        255 len-2, then the bytes
//	len 249 and up      255 247+t, then len-start in t+1 bytes little-endian,
//	                    then the bytes
//
// where t is the tier (0 to 6) and start is the shortest length in it: 249,
// 505, 66041, 16843257, 4311810553, 1103823438329 and 282578800148985.
func (s Serializer) CompactSlice(data []byte) {
	ln := len(data)
	if ln == 0 {
		s.Byte(sliceNil)
		return
	}
	if ln == 1 && data[0] < sliceNil {
		s.Byte(data[0])
		return
	}
	s.Byte(sliceExt)
	if ln == 1 {
		s.Byte(data[0])
		return
	}
	s.sliceLen(uint64(ln))
	s.Slice(data)
}

func (s Serializer) sliceLen(ln uint64) {
	if ln < tierStart[0] {
		s.Byte(byte(ln - shortMin))
		return
	}
	t := tier(ln)
	s.Byte(sliceMeta + t)
	s.Uint(t+1, ln-tierStart[t])
}

func tier(ln uint64) byte {
	t := byte(len(tierStart) - 1)
	for ln < tierStart[t] {
		t--
	}
	return t
}

// CompactSlice reads a byte slice in compact form, see
// Serializer.CompactSlice. The result shares memory with the data and does not
// allocate. An empty slice is read as nil. If the length is longer than the
// data that remains, the rest of the data is returned.
func (d Deserializer) CompactSlice() []byte {
	b := d.Byte()
	if b == sliceNil {
		return nil
	}
	if b != sliceExt {
		return d.Data[d.Idx-1 : d.Idx]
	}
	b = d.Byte()
	if b >= sliceNil {
		return d.Data[d.Idx-1 : d.Idx]
	}
	ln := d.sliceLen(b)
	if rem := uint64(len(d.Data) - d.Idx); ln > rem {
		ln = rem
	}
	return d.Slice(int(ln))
}

func (d Deserializer) sliceLen(b byte) uint64 {
	if b < sliceMeta {
		return uint64(b) + shortMin
	}
	t := b - sliceMeta
	return tierStart[t] + d.Uint(t+1)
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
	ln := uint64(len(data))
	if ln == 0 || (ln == 1 && data[0] < sliceNil) {
		return 1
	}
	if ln == 1 {
		return 2
	}
	if ln < tierStart[0] {
		return 2 + ln
	}
	return 3 + uint64(tier(ln)) + ln
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

// SizeUint64 returns the size of x in CompactUint64 form.
func SizeUint64(x uint64) uint64 {
	if x < uintPrefix {
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
