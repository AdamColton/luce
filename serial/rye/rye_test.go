package rye_test

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/adamcolton/luce/serial/rye"
	"github.com/stretchr/testify/assert"
)

func TestByte(t *testing.T) {
	b := byte(123)
	s := &rye.Serializer{
		Size: 5,
	}
	s.Make()
	s.Byte(b)
	assert.Equal(t, b, s.Data[0])
}

func TestRoundTripN(t *testing.T) {
	s := &rye.Serializer{
		Size: 8,
	}
	s.Make()

	u8 := uint8(0xaa)
	s.Uint8(u8)
	assert.Equal(t, []byte{0xaa, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, s.Data)
	s.Idx = 0

	u16 := uint16(0xabab)
	s.Uint16(u16)
	assert.Equal(t, []byte{0xab, 0xab, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, s.Data)
	s.Idx = 0

	u32 := uint32(0xabcd1234)
	s.Uint32(u32)
	assert.Equal(t, []byte{0x34, 0x12, 0xcd, 0xab, 0x0, 0x0, 0x0, 0x0}, s.Data)
	s.Idx = 0

	u64 := uint64(0xabcd1234)
	s.Uint64(u64)
	assert.Equal(t, []byte{0x34, 0x12, 0xcd, 0xab, 0x0, 0x0, 0x0, 0x0}, s.Data)
	s.Idx = 0

	i8 := int8(-32)
	s.Int8(i8)
	assert.Equal(t, []byte{0xe0, 0x12, 0xcd, 0xab, 0x0, 0x0, 0x0, 0x0}, s.Data)
	s.Idx = 0

	i16 := int16(-0x1bab)
	s.Int16(i16)
	assert.Equal(t, []byte{0x55, 0xe4, 0xcd, 0xab, 0x0, 0x0, 0x0, 0x0}, s.Data)
	s.Idx = 0

	i32 := int32(-0x2bcd1234)
	s.Int32(i32)
	assert.Equal(t, []byte{0xcc, 0xed, 0x32, 0xd4, 0x0, 0x0, 0x0, 0x0}, s.Data)
	s.Idx = 0

	i64 := int64(-0x3cd1234)
	s.Int64(i64)
	assert.Equal(t, []byte{0xcc, 0xed, 0x32, 0xfc, 0xff, 0xff, 0xff, 0xff}, s.Data)
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
	s := &rye.Serializer{
		Size: 8,
	}
	s.Make()
	clear := func() {
		s.Idx = 0
		s.Data[0] = 0
		s.Data[1] = 0
		s.Data[2] = 0
		s.Data[3] = 0
		s.Data[4] = 0
		s.Data[5] = 0
		s.Data[6] = 0
		s.Data[7] = 0
	}

	for _, tc := range tt {
		t.Run(fmt.Sprint(tc), func(t *testing.T) {
			expected := make([]byte, 8)
			rye.Serialize.Uint64(expected, tc)

			clear()
			size := s.Uint(0, tc)
			assert.Equal(t, expected, s.Data)

			clear()
			s.Uint64(tc)
			assert.Equal(t, expected, s.Data)

			if size <= 4 {
				clear()
				s.Uint32(uint32(tc))
				assert.Equal(t, expected, s.Data)
			}
			if size <= 2 {
				clear()
				s.Uint16(uint16(tc))
				assert.Equal(t, expected, s.Data)
			}
			if size <= 1 {
				clear()
				s.Uint8(uint8(tc))
				assert.Equal(t, expected, s.Data)
			}
		})
	}
}

// func TestCompactUint64(t *testing.T) {
// 	tt := []uint64{
// 		0,
// 		100,
// 		128,
// 		129,
// 		0xaa, 0xff, 0x100,
// 		0xaaa, 0xfff, 0x1000,
// 		0xaaaa, 0xffff, 0x10000,
// 		0xaaaaa, 0xfffff, 0x100000,
// 		0xaaaaaa, 0xffffff, 0x1000000,
// 		0xaaaaaaa, 0xfffffff, 0x10000000,
// 		0xaaaaaaaa, 0xffffffff, 0x100000000,
// 		0xaaaaaaaaa, 0xfffffffff, 0x1000000000,
// 		0xaaaaaaaaaa, 0xffffffffff, 0x10000000000,
// 		0xaaaaaaaaaaa, 0xfffffffffff, 0x100000000000,
// 		0xaaaaaaaaaaaa, 0xffffffffffff, 0x1000000000000,
// 		0xaaaaaaaaaaaaa, 0xfffffffffffff, 0x10000000000000,
// 		0xaaaaaaaaaaaaaa, 0xffffffffffffff, 0x100000000000000,
// 		0xaaaaaaaaaaaaaaa, 0xfffffffffffffff, 0x1000000000000000,
// 		0xaaaaaaaaaaaaaaaa, 0xffffffffffffffff,
// 	}

// 	s := &rye.Serializer{
// 		Size: 9,
// 	}
// 	s.Make()

// 	for _, tc := range tt {
// 		t.Run(fmt.Sprint(tc), func(t *testing.T) {
// 			s.Idx = 0
// 			s.CompactUint64(tc)
// 			assert.Equal(t, uint64(s.Idx), Compact.SizeUint64(tc))
// 			assert.Equal(t, tc, NewDeserializer(s.Data).CompactUint64())

// 			// Check that CompactUint64 and CompactSlice logic is compatable
// 			b := NewDeserializer(s.Data).CompactSlice()
// 			assert.Equal(t, tc, NewDeserializer(b).Uint(byte(len(b))))
// 		})
// 	}
// }

// func TestCompactSlice(t *testing.T) {
// 	tt := [][]byte{
// 		nil,
// 		[]byte("A"),
// 		[]byte("z"),
// 		[]byte("aa"),
// 		[]byte("abc"),
// 		[]byte("this is a longer test"),
// 	}

// 	long := ""
// 	for len(long) < 120 {
// 		long += "test "
// 	}
// 	assert.Len(t, long, 120)
// 	tt = append(tt, []byte(long))
// 	tt = append(tt, []byte(long+"!"))
// 	tt = append(tt, []byte(long+"!!"))
// 	for len(long) < 256 {
// 		long += "test "
// 	}
// 	tt = append(tt, []byte(long))

// 	for _, tc := range tt {
// 		t.Run(string(tc), func(t *testing.T) {
// 			s := &rye.Serializer{
// 				Size: len(tc) + 3,
// 			}
// 			s.Make()
// 			s.CompactSlice(tc)
// 			assert.Equal(t, uint64(s.Idx), Compact.Size(tc))
// 			assert.Equal(t, tc, NewDeserializer(s.Data).CompactSlice())

// 			s.Data = nil
// 			s.Idx = 0
// 			s.Make()
// 			str := string(tc)
// 			s.CompactString(str)
// 			assert.Equal(t, str, NewDeserializer(s.Data).CompactString())
// 		})
// 	}
// }

func TestFuzzFloats(t *testing.T) {
	s := &rye.Serializer{
		Size: 8,
	}
	s.Make()
	for i := 0; i < 100; i++ {
		f32 := rand.Float32()
		s.Idx = 0
		s.Float32(f32)
		u32 := rye.Deserialize.Uint32(s.Data)
		got32 := math.Float32frombits(u32)
		assert.Equal(t, f32, got32)

		f64 := rand.Float64()
		s.Idx = 0
		s.Float64(f64)
		u64 := rye.Deserialize.Uint64(s.Data)
		got64 := math.Float64frombits(u64)
		assert.Equal(t, f64, got64)
	}
}

// func TestCompactInt64(t *testing.T) {
// 	tt := []int64{
// 		0,
// 		100, 128, 129,
// 		-100, -128, -129,
// 		0xaa, 0xff, 0x100,
// 		0xaaa, 0xfff, 0x1000,
// 		0xaaaa, 0xffff, 0x10000,
// 		0xaaaaa, 0xfffff, 0x100000,
// 		0xaaaaaa, 0xffffff, 0x1000000,
// 		0xaaaaaaa, 0xfffffff, 0x10000000,
// 		0xaaaaaaaa, 0xffffffff, 0x100000000,
// 		0xaaaaaaaaa, 0xfffffffff, 0x1000000000,
// 		0xaaaaaaaaaa, 0xffffffffff, 0x10000000000,
// 		0xaaaaaaaaaaa, 0xfffffffffff, 0x100000000000,
// 		0xaaaaaaaaaaaa, 0xffffffffffff, 0x1000000000000,
// 		0xaaaaaaaaaaaaa, 0xfffffffffffff, 0x10000000000000,
// 		0xaaaaaaaaaaaaaa, 0xffffffffffffff, 0x100000000000000,
// 		0xaaaaaaaaaaaaaaa, 0xfffffffffffffff, 0x1000000000000000,
// 		-0xaa, -0xff, -0x100,
// 		-0xaaa, -0xfff, -0x1000,
// 		-0xaaaa, -0xffff, -0x10000,
// 		-0xaaaaa, -0xfffff, -0x100000,
// 		-0xaaaaaa, -0xffffff, -0x1000000,
// 		-0xaaaaaaa, -0xfffffff, -0x10000000,
// 		-0xaaaaaaaa, -0xffffffff, -0x100000000,
// 		-0xaaaaaaaaa, -0xfffffffff, -0x1000000000,
// 		-0xaaaaaaaaaa, -0xffffffffff, -0x10000000000,
// 		-0xaaaaaaaaaaa, -0xfffffffffff, -0x100000000000,
// 		-0xaaaaaaaaaaaa, -0xffffffffffff, -0x1000000000000,
// 		-0xaaaaaaaaaaaaa, -0xfffffffffffff, -0x10000000000000,
// 		-0xaaaaaaaaaaaaaa, -0xffffffffffffff, -0x100000000000000,
// 		-0xaaaaaaaaaaaaaaa, -0xfffffffffffffff, -0x1000000000000000,
// 	}

// 	s := &rye.Serializer{
// 		Size: 9,
// 	}
// 	s.Make()

// 	for _, tc := range tt {
// 		t.Run(fmt.Sprint(tc), func(t *testing.T) {
// 			s.Idx = 0
// 			s.CompactInt64(tc)
// 			assert.Equal(t, uint64(s.Idx), Compact.SizeInt64(tc))
// 			assert.Equal(t, tc, NewDeserializer(s.Data).CompactInt64())
// 		})
// 	}
// }

// func TestCompactSub(t *testing.T) {
// 	a := make([]byte, 10)
// 	b := make([]byte, 10)
// 	rand.Read(a)
// 	rand.Read(b)

// 	c := make([]uint64, 10)
// 	inner := &rye.Serializer{}
// 	for i := range c {
// 		c[i] = rand.Uint64()
// 		inner.Size += int(Compact.SizeUint64(c[i]))
// 	}
// 	inner.Make()
// 	for _, u := range c {
// 		inner.CompactUint64(u)
// 	}

// 	s := &rye.Serializer{
// 		Size: int(Compact.Size(a) + Compact.Size(b) + Compact.Size(inner.Data)),
// 	}
// 	s.Make()
// 	s.CompactSlice(a)
// 	s.CompactSlice(inner.Data)
// 	s.CompactSlice(b)

// 	d := NewDeserializer(s.Data)

// 	assert.Equal(t, a, d.CompactSlice())
// 	sub := d.CompactSub()
// 	assert.Equal(t, b, d.CompactSlice())
// 	for i := 0; !sub.Done(); i++ {
// 		assert.Equal(t, c[i], sub.CompactUint64())
// 	}
// }

// func TestDeserializerSub(t *testing.T) {
// 	d := NewDeserializer([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
// 	s := d.Sub(5)
// 	assert.Equal(t, byte(6), d.Byte())
// 	assert.Equal(t, byte(1), s.Byte())
// 	assert.Equal(t, byte(2), s.Byte())
// 	assert.Equal(t, byte(7), d.Byte())
// }

func TestSerializerSub(t *testing.T) {
	s := (&rye.Serializer{Size: 10}).Make()
	sub := s.Sub(5)
	for i := byte(0); i < 5; i++ {
		s.Byte(i + 6)
		sub.Byte(i + 1)
	}
	assert.Equal(t, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, s.Data)
}

// func TestDeserializerString(t *testing.T) {
// 	d := NewDeserializer([]byte("Hello, World"))
// 	assert.Equal(t, "Hello", d.String(5))
// 	assert.Equal(t, byte(','), d.Byte())
// }

func TestSerializerString(t *testing.T) {
	str := "Hello, World"
	s := (&rye.Serializer{Size: len(str)}).Make()
	s.String(str)
	assert.Equal(t, len(str), s.Idx)
	assert.Equal(t, str, string(s.Data))
}

// func TestSerializerCheckFree(t *testing.T) {
// 	s := (&rye.Serializer{Size: 5}).Make()
// 	assert.Len(t, s.Data, 5)
// 	s.CheckFree(4)
// 	assert.Len(t, s.Data, 5)
// 	s.CheckFree(10)
// 	assert.Len(t, s.Data, 10)
// }
