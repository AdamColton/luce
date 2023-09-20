package compact_test

import (
	crand "crypto/rand"
	"fmt"
	"math/rand"
	"testing"

	"github.com/adamcolton/luce/serial/rye"
	"github.com/adamcolton/luce/serial/rye/compact"
	"github.com/stretchr/testify/assert"
)

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

	s := compact.NewSerializer(9)
	s.Make()

	for _, tc := range tt {
		t.Run(fmt.Sprint(tc), func(t *testing.T) {
			s.Idx = 0
			s.CompactUint64(tc)
			assert.Equal(t, uint64(s.Idx), compact.SizeUint64(tc))
			assert.Equal(t, tc, compact.NewDeserializer(s.Data).CompactUint64())

			// Check that CompactUint64 and CompactSlice logic is compatable
			b := compact.NewDeserializer(s.Data).CompactSlice()
			assert.Equal(t, tc, rye.NewDeserializer(b).Uint(byte(len(b))))
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
			s := compact.NewSerializer(len(tc) + 3)
			s.Make()
			s.CompactSlice(tc)
			assert.Equal(t, uint64(s.Idx), compact.Size(tc))
			assert.Equal(t, tc, compact.NewDeserializer(s.Data).CompactSlice())

			s.Data = nil
			s.Idx = 0
			s.Make()
			str := string(tc)
			s.CompactString(str)
			assert.Equal(t, str, compact.NewDeserializer(s.Data).CompactString())
		})
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

	s := compact.NewSerializer(9)
	s.Make()

	for _, tc := range tt {
		t.Run(fmt.Sprint(tc), func(t *testing.T) {
			s.Idx = 0
			s.CompactInt64(tc)
			assert.Equal(t, uint64(s.Idx), compact.SizeInt64(tc))
			assert.Equal(t, tc, compact.NewDeserializer(s.Data).CompactInt64())
		})
	}
}

func TestCompactSub(t *testing.T) {
	a := make([]byte, 10)
	b := make([]byte, 10)
	crand.Read(a)
	crand.Read(b)

	c := make([]uint64, 10)
	inner := compact.NewSerializer(0)
	for i := range c {
		c[i] = rand.Uint64()
		inner.Size += int(compact.SizeUint64(c[i]))
	}
	inner.Make()
	for _, u := range c {
		inner.CompactUint64(u)
	}

	s := compact.NewSerializer(int(compact.Size(a) + compact.Size(b) + compact.Size(inner.Data)))
	s.Make()
	s.CompactSlice(a)
	s.CompactSlice(inner.Data)
	s.CompactSlice(b)

	d := compact.NewDeserializer(s.Data)

	assert.Equal(t, a, d.CompactSlice())
	sub := d.CompactSub()
	assert.Equal(t, b, d.CompactSlice())
	for i := 0; !sub.Done(); i++ {
		assert.Equal(t, c[i], sub.CompactUint64())
	}
}

func TestSizeUint0(t *testing.T) {
	assert.Equal(t, uint64(1), compact.SizeUint(0))
}
