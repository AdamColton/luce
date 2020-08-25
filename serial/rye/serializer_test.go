package rye

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundTripByte(t *testing.T) {
	b := byte(123)
	s := &Serializer{
		Size: 5,
	}
	s.Make()
	s.Byte(b)
	assert.Equal(t, b, s.Data[0])

	d := NewDeserializer(s.Data)
	assert.Equal(t, b, d.Byte())
}

func TestRoundTripN(t *testing.T) {
	s := &Serializer{
		Size: 8,
	}
	s.Make()

	u8 := uint8(0xaa)
	s.Uint8(u8)
	assert.Equal(t, u8, NewDeserializer(s.Data).Uint8())
	s.Idx = 0

	u16 := uint16(0xabab)
	s.Uint16(u16)
	assert.Equal(t, u16, NewDeserializer(s.Data).Uint16())
	s.Idx = 0

	u32 := uint32(0xabcd1234)
	s.Uint32(u32)
	assert.Equal(t, u32, NewDeserializer(s.Data).Uint32())
	s.Idx = 0

	u64 := uint64(0xabcd1234)
	s.Uint64(u64)
	assert.Equal(t, u64, NewDeserializer(s.Data).Uint64())
	s.Idx = 0

	i8 := int8(-32)
	s.Int8(i8)
	assert.Equal(t, i8, NewDeserializer(s.Data).Int8())
	s.Idx = 0

	i16 := int16(-0x1bab)
	s.Int16(i16)
	assert.Equal(t, i16, NewDeserializer(s.Data).Int16())
	s.Idx = 0

	i32 := int32(-0x2bcd1234)
	s.Int32(i32)
	assert.Equal(t, i32, NewDeserializer(s.Data).Int32())
	s.Idx = 0

	i64 := int64(-0x3cd1234)
	s.Int64(i64)
	assert.Equal(t, i64, NewDeserializer(s.Data).Int64())
	s.Idx = 0
}

func TestRoundTripUint(t *testing.T) {
	tt := []uint64{
		0,
		100,
		128,
		129,
		0xaa, 0xff, 0x100,
		0xaaa, 0xfff, 0x1000,
		0xaaaa, 0xffff, 0x10000,
		0xaaaaa, 0xfffff, 0x100000,
		0xaaaaaa, 0xffffff, 0x1000000,
		0xaaaaaaa, 0xfffffff, 0x10000000,
		0xaaaaaaaa, 0xffffffff, 0x100000000,
		0xaaaaaaaaa, 0xfffffffff, 0x1000000000,
		0xaaaaaaaaaa, 0xffffffffff, 0x10000000000,
		0xaaaaaaaaaaa, 0xfffffffffff, 0x100000000000,
		0xaaaaaaaaaaaa, 0xffffffffffff, 0x1000000000000,
		0xaaaaaaaaaaaaa, 0xfffffffffffff, 0x10000000000000,
		0xaaaaaaaaaaaaaa, 0xffffffffffffff, 0x100000000000000,
		0xaaaaaaaaaaaaaaa, 0xfffffffffffffff, 0x1000000000000000,
		0xaaaaaaaaaaaaaaaa, 0xffffffffffffffff,
	}
	s := &Serializer{
		Size: 10,
	}
	s.Make()
	end := byte(31)

	for _, tc := range tt {
		t.Run(fmt.Sprint(tc), func(t *testing.T) {
			s.Idx = 0
			size := s.Uint(0, tc)
			assert.Equal(t, uint64(size), SizeUint(tc))
			s.Byte(end)
			d := NewDeserializer(s.Data)
			assert.Equal(t, tc, d.Uint(size))
			assert.Equal(t, end, d.Byte())
			for size++; size < 8; size++ {
				s.Idx = 0
				s.Uint(size, tc)
				assert.Equal(t, tc, NewDeserializer(s.Data).Uint(size))
			}
		})
	}
}

func TestCompactUint64(t *testing.T) {
	tt := []uint64{
		0,
		100,
		128,
		129,
		0xaa, 0xff, 0x100,
		0xaaa, 0xfff, 0x1000,
		0xaaaa, 0xffff, 0x10000,
		0xaaaaa, 0xfffff, 0x100000,
		0xaaaaaa, 0xffffff, 0x1000000,
		0xaaaaaaa, 0xfffffff, 0x10000000,
		0xaaaaaaaa, 0xffffffff, 0x100000000,
		0xaaaaaaaaa, 0xfffffffff, 0x1000000000,
		0xaaaaaaaaaa, 0xffffffffff, 0x10000000000,
		0xaaaaaaaaaaa, 0xfffffffffff, 0x100000000000,
		0xaaaaaaaaaaaa, 0xffffffffffff, 0x1000000000000,
		0xaaaaaaaaaaaaa, 0xfffffffffffff, 0x10000000000000,
		0xaaaaaaaaaaaaaa, 0xffffffffffffff, 0x100000000000000,
		0xaaaaaaaaaaaaaaa, 0xfffffffffffffff, 0x1000000000000000,
		0xaaaaaaaaaaaaaaaa, 0xffffffffffffffff,
	}

	s := &Serializer{
		Size: 9,
	}
	s.Make()

	for _, tc := range tt {
		t.Run(fmt.Sprint(tc), func(t *testing.T) {
			s.Idx = 0
			s.CompactUint64(tc)
			assert.Equal(t, uint64(s.Idx), SizeCompactUint64(tc))
			assert.Equal(t, tc, NewDeserializer(s.Data).CompactUint64())

			// Check that CompactUint64 and CompactSlice logic is compatable
			b := NewDeserializer(s.Data).CompactSlice()
			assert.Equal(t, tc, NewDeserializer(b).Uint(byte(len(b))))
		})
	}
}

func TestCompactSlice(t *testing.T) {
	tt := [][]byte{
		nil,
		[]byte("A"),
		[]byte("z"),
		[]byte("aa"),
		[]byte("abc"),
		[]byte("this is a longer test"),
	}

	long := ""
	for len(long) < 120 {
		long += "test "
	}
	assert.Len(t, long, 120)
	tt = append(tt, []byte(long))
	tt = append(tt, []byte(long+"!"))
	tt = append(tt, []byte(long+"!!"))
	for len(long) < 256 {
		long += "test "
	}
	tt = append(tt, []byte(long))

	for _, tc := range tt {
		t.Run(string(tc), func(t *testing.T) {
			s := &Serializer{
				Size: len(tc) + 3,
			}
			s.Make()
			s.CompactSlice(tc)
			assert.Equal(t, uint64(s.Idx), Size(tc))
			assert.Equal(t, tc, NewDeserializer(s.Data).CompactSlice())

			s.Data = nil
			s.Idx = 0
			s.Make()
			str := string(tc)
			s.CompactString(str)
			assert.Equal(t, str, NewDeserializer(s.Data).CompactString())
		})
	}
}

func TestFuzzFloats(t *testing.T) {
	s := &Serializer{
		Size: 8,
	}
	s.Make()
	d := NewDeserializer(s.Data)
	for i := 0; i < 100; i++ {
		f32 := rand.Float32()
		s.Idx = 0
		s.Float32(f32)
		d.Idx = 0
		assert.Equal(t, f32, d.Float32())

		f64 := rand.Float64()
		s.Idx = 0
		s.Float64(f64)
		d.Idx = 0
		assert.Equal(t, f64, d.Float64())
	}
}

func TestCompactInt64(t *testing.T) {
	tt := []int64{
		0,
		100, 128, 129,
		-100, -128, -129,
		0xaa, 0xff, 0x100,
		0xaaa, 0xfff, 0x1000,
		0xaaaa, 0xffff, 0x10000,
		0xaaaaa, 0xfffff, 0x100000,
		0xaaaaaa, 0xffffff, 0x1000000,
		0xaaaaaaa, 0xfffffff, 0x10000000,
		0xaaaaaaaa, 0xffffffff, 0x100000000,
		0xaaaaaaaaa, 0xfffffffff, 0x1000000000,
		0xaaaaaaaaaa, 0xffffffffff, 0x10000000000,
		0xaaaaaaaaaaa, 0xfffffffffff, 0x100000000000,
		0xaaaaaaaaaaaa, 0xffffffffffff, 0x1000000000000,
		0xaaaaaaaaaaaaa, 0xfffffffffffff, 0x10000000000000,
		0xaaaaaaaaaaaaaa, 0xffffffffffffff, 0x100000000000000,
		0xaaaaaaaaaaaaaaa, 0xfffffffffffffff, 0x1000000000000000,
		-0xaa, -0xff, -0x100,
		-0xaaa, -0xfff, -0x1000,
		-0xaaaa, -0xffff, -0x10000,
		-0xaaaaa, -0xfffff, -0x100000,
		-0xaaaaaa, -0xffffff, -0x1000000,
		-0xaaaaaaa, -0xfffffff, -0x10000000,
		-0xaaaaaaaa, -0xffffffff, -0x100000000,
		-0xaaaaaaaaa, -0xfffffffff, -0x1000000000,
		-0xaaaaaaaaaa, -0xffffffffff, -0x10000000000,
		-0xaaaaaaaaaaa, -0xfffffffffff, -0x100000000000,
		-0xaaaaaaaaaaaa, -0xffffffffffff, -0x1000000000000,
		-0xaaaaaaaaaaaaa, -0xfffffffffffff, -0x10000000000000,
		-0xaaaaaaaaaaaaaa, -0xffffffffffffff, -0x100000000000000,
		-0xaaaaaaaaaaaaaaa, -0xfffffffffffffff, -0x1000000000000000,
	}

	s := &Serializer{
		Size: 9,
	}
	s.Make()

	for _, tc := range tt {
		t.Run(fmt.Sprint(tc), func(t *testing.T) {
			s.Idx = 0
			s.CompactInt64(tc)
			assert.Equal(t, uint64(s.Idx), SizeCompactInt64(tc))
			assert.Equal(t, tc, NewDeserializer(s.Data).CompactInt64())
		})
	}
}

func TestCompactSub(t *testing.T) {
	a := make([]byte, 10)
	b := make([]byte, 10)
	rand.Read(a)
	rand.Read(b)

	c := make([]uint64, 10)
	inner := &Serializer{}
	for i := range c {
		c[i] = rand.Uint64()
		inner.Size += int(SizeCompactUint64(c[i]))
	}
	inner.Make()
	for _, u := range c {
		inner.CompactUint64(u)
	}

	s := &Serializer{
		Size: int(Size(a) + Size(b) + Size(inner.Data)),
	}
	s.Make()
	s.CompactSlice(a)
	s.CompactSlice(inner.Data)
	s.CompactSlice(b)

	d := NewDeserializer(s.Data)

	assert.Equal(t, a, d.CompactSlice())
	sub := d.CompactSub()
	assert.Equal(t, b, d.CompactSlice())
	for i := 0; !sub.Done(); i++ {
		assert.Equal(t, c[i], sub.CompactUint64())
	}
}
