package merkle_test

import (
	"crypto/rand"
	"crypto/sha256"
	"io"
	mrand "math/rand"
	"testing"

	"github.com/adamcolton/luce/ds/merkle"
	"github.com/adamcolton/luce/lerr"
	"github.com/stretchr/testify/assert"
)

func TestLeaf(t *testing.T) {
	var maxLeafSize uint32 = 200
	var ln uint32 = 5000
	data := make([]byte, ln)
	rand.Read(data)

	b := merkle.NewBuilder(maxLeafSize, sha256.New)
	m := b.Build(nil)
	assert.Nil(t, m)
	m = b.Build(data)
	assert.Equal(t, data, m.Data())

	start := 0
	h := sha256.New()
	buf := make([]byte, 32)
	leaves := m.Leaves()
	for i := 0; i < leaves; i++ {
		l := m.Leaf(i)
		d := l.Data
		end := start + len(d)
		assert.Equal(t, i, int(l.Index))
		assert.Equal(t, data[start:end], d)
		assert.Equal(t, m.Digest(), l.Digest(h, buf))
		start = end
	}
	assert.Nil(t, m.Leaf(leaves))
}

func TestAssembler(t *testing.T) {
	var maxLeafSize uint32 = 200
	var ln uint32 = 5000
	data := make([]byte, ln)
	rand.Read(data)

	m := merkle.NewBuilder(maxLeafSize, sha256.New).Build(data)
	assert.Equal(t, data, m.Data())

	leaves := make([]*merkle.Leaf, m.Leaves())
	for i := range leaves {
		leaves[i] = m.Leaf(i)
	}
	mrand.New(mrand.NewSource(31415)).Shuffle(len(leaves), func(i, j int) {
		leaves[i], leaves[j] = leaves[j], leaves[i]
	})

	a := m.Description().Assembler(sha256.New())

	l := m.Leaf(0)
	l.Index = uint32(m.Leaves()) + 1
	assert.False(t, a.Add(l))

	for _, l := range leaves {
		done, tr := a.Done()
		assert.False(t, done)
		assert.Nil(t, tr)
		assert.True(t, a.Add(l), l.Index)
		assert.False(t, a.Add(l))
	}
	done, tr := a.Done()
	assert.True(t, done)
	assert.Equal(t, data, tr.Data())
}

func TestReaderSeeker(t *testing.T) {
	var maxLeafSize uint32 = 200
	var ln uint32 = 5000
	data := make([]byte, ln)
	rand.Read(data)

	m := merkle.NewBuilder(maxLeafSize, sha256.New).Build(data)
	got, err := io.ReadAll(m)
	assert.NoError(t, err)
	assert.Equal(t, data, got)

	x := int(ln / 10)
	tt := []struct {
		name       string
		whence     int
		offset     int64
		bufLen     int
		start, end int
	}{
		{
			name:   "SeekStart",
			whence: io.SeekStart,
			offset: int64(x),
			bufLen: x,
			start:  x,
			end:    2 * x,
		},
		{
			name:   "SeekCurrent",
			whence: io.SeekCurrent,
			offset: int64(x),
			bufLen: x,
			start:  3 * x,
			end:    4 * x,
		},
		{
			name:   "SeekEnd",
			whence: io.SeekEnd,
			offset: -int64(2 * x),
			bufLen: x,
			start:  int(ln) - 2*x,
			end:    int(ln) - x,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			pos, err := m.Seek(tc.offset, tc.whence)
			assert.NoError(t, err)
			assert.Equal(t, int64(tc.start), pos)
			buf := make([]byte, tc.bufLen)
			n, err := m.Read(buf)
			assert.NoError(t, err)
			assert.Equal(t, tc.end-tc.start, n)
			assert.Equal(t, data[tc.start:tc.end], buf)
		})
	}

	n, _ := m.Seek(1, io.SeekEnd)
	assert.Equal(t, int64(ln), n)
	n, _ = m.Seek(-1, io.SeekStart)
	assert.Equal(t, int64(0), n)
	_, err = m.Seek(100, 10)
	assert.Equal(t, lerr.ErrBadWhence(10), err)

}
