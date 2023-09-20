package rye_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/adamcolton/luce/serial/rye"
	"github.com/stretchr/testify/assert"
)

func TestRoundTripByte(t *testing.T) {
	b := byte(123)
	s := &rye.Serializer{
		Size: 5,
	}
	s.Make()
	s.Byte(b)
	assert.Equal(t, b, s.Data[0])

	d := rye.NewDeserializer(s.Data)
	assert.Equal(t, b, d.Byte())
}

func TestRoundTripN(t *testing.T) {
	s := &rye.Serializer{
		Size: 8,
	}
	s.Make()

	u8 := uint8(0xaa)
	s.Uint8(u8)
	assert.Equal(t, u8, rye.NewDeserializer(s.Data).Uint8())
	s.Idx = 0

	u16 := uint16(0xabab)
	s.Uint16(u16)
	assert.Equal(t, u16, rye.NewDeserializer(s.Data).Uint16())
	s.Idx = 0

	u32 := uint32(0xabcd1234)
	s.Uint32(u32)
	assert.Equal(t, u32, rye.NewDeserializer(s.Data).Uint32())
	s.Idx = 0

	u64 := uint64(0xabcd1234)
	s.Uint64(u64)
	assert.Equal(t, u64, rye.NewDeserializer(s.Data).Uint64())
	s.Idx = 0

	i8 := int8(-32)
	s.Int8(i8)
	assert.Equal(t, i8, rye.NewDeserializer(s.Data).Int8())
	s.Idx = 0

	i16 := int16(-0x1bab)
	s.Int16(i16)
	assert.Equal(t, i16, rye.NewDeserializer(s.Data).Int16())
	s.Idx = 0

	i32 := int32(-0x2bcd1234)
	s.Int32(i32)
	assert.Equal(t, i32, rye.NewDeserializer(s.Data).Int32())
	s.Idx = 0

	i64 := int64(-0x3cd1234)
	s.Int64(i64)
	assert.Equal(t, i64, rye.NewDeserializer(s.Data).Int64())
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
		Size: 10,
	}
	s.Make()
	end := byte(31)

	for _, tc := range tt {
		t.Run(fmt.Sprint(tc), func(t *testing.T) {
			s.Idx = 0
			size := s.Uint(0, tc)
			s.Byte(end)
			d := rye.NewDeserializer(s.Data)
			assert.Equal(t, tc, d.Uint(size))
			assert.Equal(t, end, d.Byte())
			for size++; size < 8; size++ {
				s.Idx = 0
				s.Uint(size, tc)
				assert.Equal(t, tc, rye.NewDeserializer(s.Data).Uint(size))
			}
		})
	}
}

func TestFuzzFloats(t *testing.T) {
	s := &rye.Serializer{
		Size: 8,
	}
	s.Make()
	d := rye.NewDeserializer(s.Data)
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

func TestDeserializerSub(t *testing.T) {
	d := rye.NewDeserializer([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	s := d.Sub(5)
	assert.Equal(t, byte(6), d.Byte())
	assert.Equal(t, byte(1), s.Byte())
	assert.Equal(t, byte(2), s.Byte())
	assert.Equal(t, byte(7), d.Byte())
}

func TestSerializerSub(t *testing.T) {
	s := (&rye.Serializer{Size: 10}).Make()
	sub := s.Sub(5)
	for i := byte(0); i < 5; i++ {
		s.Byte(i + 6)
		sub.Byte(i + 1)
	}
	assert.Equal(t, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, s.Data)
}

func TestDeserializerString(t *testing.T) {
	d := rye.NewDeserializer([]byte("Hello, World"))
	assert.Equal(t, "Hello", d.String(5))
	assert.Equal(t, byte(','), d.Byte())
}

func TestSerializerString(t *testing.T) {
	str := "Hello, World"
	s := (&rye.Serializer{Size: len(str)}).Make()
	s.String(str)
	assert.Equal(t, len(str), s.Idx)
	assert.Equal(t, str, string(s.Data))
}

func TestDeserializerDone(t *testing.T) {
	in := rye.NewDeserializer([]byte("testing"))
	out := make([]byte, 0, 7)
	for !in.Done() {
		out = append(out, in.Byte())
	}
	assert.Equal(t, "testing", string(out))
}
