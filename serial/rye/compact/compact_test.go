package compact_test

import (
	"bytes"
	crand "crypto/rand"
	"fmt"
	"math/rand"
	"testing"

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
		})
	}
}

func TestCompactSlice(t *testing.T) {
	tt := [][]byte{
		nil,
		[]byte("A"),
		[]byte("z"),
		{253}, {254}, {255},
		[]byte("aa"),
		[]byte("abc"),
		[]byte("this is a longer test"),
	}
	// Lengths on either side of each boundary between header sizes.
	for _, ln := range []int{248, 249, 504, 505, 66040, 66041} {
		tt = append(tt, bytes.Repeat([]byte("test "), ln/5+1)[:ln])
	}

	for _, tc := range tt {
		t.Run(fmt.Sprint(len(tc)), func(t *testing.T) {
			s := compact.NewSerializer(int(compact.Size(tc)))
			s.Make()
			s.CompactSlice(tc)
			assert.Equal(t, len(s.Data), s.Idx)
			assert.True(t, bytes.Equal(tc, compact.NewDeserializer(s.Data).CompactSlice()))

			str := string(tc)
			s = compact.NewSerializer(int(compact.Size([]byte(str))))
			s.Make()
			s.CompactString(str)
			assert.True(t, str == compact.NewDeserializer(s.Data).CompactString())
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

// A decoded slice points into the data it was read from, so reading never
// allocates, even for a single byte.
func TestCompactSliceShares(t *testing.T) {
	data := []byte{7, 255, 255, 255, 0, 'h', 'i'}
	d := compact.NewDeserializer(data)

	var seven, ff, hi []byte
	allocs := testing.AllocsPerRun(10, func() {
		d.Idx = 0
		seven = d.CompactSlice()
		ff = d.CompactSlice()
		hi = d.CompactSlice()
	})
	assert.Zero(t, allocs)
	assert.Equal(t, []byte{7}, seven)
	assert.Same(t, &data[0], &seven[0])
	assert.Equal(t, []byte{255}, ff)
	assert.Same(t, &data[2], &ff[0])
	assert.Equal(t, []byte("hi"), hi)
	assert.Same(t, &data[5], &hi[0])
}

// When the data ends before the length it claims, the rest of the data is
// returned. A Sub never reads past its own data.
func TestCompactSliceTruncated(t *testing.T) {
	full := encodeSlice([]byte("hello world"))
	d := compact.NewDeserializer(full[:7])
	assert.Equal(t, []byte("hello"), d.CompactSlice())
	assert.True(t, d.Done())

	s := compact.NewSerializer(int(compact.Size(full[:7]) + compact.Size([]byte("next"))))
	s.Make()
	s.CompactSlice(full[:7])
	s.CompactString("next")
	d = compact.NewDeserializer(s.Data)
	sub := d.CompactSub()
	assert.Equal(t, []byte("hello"), sub.CompactSlice())
	assert.Equal(t, "next", d.CompactString())
}

func TestSizeUint0(t *testing.T) {
	assert.Equal(t, uint64(1), compact.SizeUint(0))
}
