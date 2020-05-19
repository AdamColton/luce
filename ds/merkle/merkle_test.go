package merkle

import (
	"crypto/sha256"
	"io"
	"math/rand"
	"sort"
	"testing"

	"github.com/testify/assert"
)

func TestByteHasher(t *testing.T) {
	b := newDataLeaf([]byte("this is a test"), 0, sha256.New())
	expected := []byte{89, 249, 159, 51, 240, 220, 187, 52, 220, 54, 194, 151,
		130, 196, 153, 59, 108, 250, 25, 202, 230, 85, 152, 51, 35, 103, 207,
		119, 148, 46, 142, 154}
	assert.Equal(t, expected, b.Digest())
}

func TestBuild(t *testing.T) {
	data := make([]byte, 50)
	rand.Read(data)

	b := NewBuilder(6, 3, sha256.New())
	n := b.Build(data)

	assert.Equal(t, data, n.Data())
}

func TestEnd2EndFuzz(t *testing.T) {
	for round := 0; round < 100; round++ {
		data := make([]byte, rand.Intn(500)+500)
		rand.Read(data)

		b := NewBuilder(uint32(rand.Intn(30)+2), byte(rand.Intn(10)+2), sha256.New())
		t1 := b.Build(data)
		td := t1.Digest()

		assert.Equal(t, data, t1.Data())

		c := t1.Count()
		lst := randList(c)

		l := t1.Leaf(lst[0])
		assert.Equal(t, td, l.Digest(sha256.New()))
		a := t1.Description().Assembler(sha256.New())

		l.Data[0] ^= 1
		assert.False(t, a.Add(l))
		l.Data[0] ^= 1

		r := l.Rows
		l.Rows = l.Rows[:len(l.Rows)-1]
		assert.False(t, a.Add(l))
		l.Rows = r

		done, t2 := a.Done()
		assert.False(t, done)
		for i, idx := range lst {
			done, t2 = a.Done()
			assert.False(t, done)
			l = t1.Leaf(idx)
			assert.True(t, a.Add(l))
			assert.Equal(t, td, l.Digest(sha256.New()))

			assert.Equal(t, expectNeed(uint32(i), lst), a.Need())

			if i == 0 {
				// test incomplete tree calls
				incpTr := Tree(a.root)
				assert.Equal(t, -1, incpTr.size())
				assert.Equal(t, maxUint32, incpTr.Count())
				assert.Equal(t, -1, incpTr.Depth())
				assert.Nil(t, incpTr.Data())
			}
		}
		done, t2 = a.Done()
		assert.True(t, done)

		// Test that re-adding a Leaf correctly checks it's validity
		for _, i := range randList(c) {
			l = t1.Leaf(i)

			d := rand.Intn(len(l.Data))
			l.Data[d]++
			assert.False(t, a.Add(l))

			l.Data[d]--
			assert.True(t, a.Add(l))

			r := l.Rows[rand.Intn(len(l.Rows))]
			dIdx := rand.Intn(len(r.Digests))
			dig := r.Digests[dIdx]
			cp := make([]byte, len(dig))
			copy(cp, dig)
			d = rand.Intn(len(dig))
			cp[d]++
			r.Digests[dIdx] = cp
			assert.False(t, a.Add(l))
		}

		assert.Equal(t, t1.Data(), t2.Data())
		assert.Equal(t, data, t1.Data())
		assert.Equal(t, data, t2.Data())
		assert.Equal(t, t1.Count(), t2.Count())
		assert.Equal(t, t1.Depth(), t2.Depth())

		buf := make([]byte, 100)
		ln, err := t1.Read(buf)
		for i := 0; err == nil; i += len(buf) {
			if err != nil {
				break
			}
			a := data[i : i+ln]
			b := buf[:ln]
			assert.Equal(t, a, b)
			ln, err = t1.Read(buf)
		}
		assert.Equal(t, err, io.EOF)

		p, err := t1.Seek(50, io.SeekStart)
		assert.NoError(t, err)
		assert.Equal(t, int64(50), p)
		_, err = t1.Read(buf)
		assert.NoError(t, err)
		assert.Equal(t, data[50:150], buf)

		p, err = t1.Seek(75, io.SeekCurrent)
		assert.NoError(t, err)
		assert.Equal(t, int64(225), p)
		_, err = t1.Read(buf)
		assert.NoError(t, err)
		assert.Equal(t, data[225:325], buf)

		p, err = t1.Seek(-90, io.SeekEnd)
		assert.NoError(t, err)
		assert.Equal(t, int64(len(data))-90, p)
		_, err = t1.Read(buf)
		assert.NoError(t, err)
		assert.Equal(t, data[p:], buf[:90])

		p, err = t1.Seek(-10, io.SeekStart)
		assert.Equal(t, int64(0), p)
		p, err = t1.Seek(1, io.SeekEnd)
		assert.Equal(t, int64(len(data)), p)

		_, err = t1.Seek(-10, 5)
		assert.Equal(t, err, ErrBadWhence)
	}
}

func randList(c uint32) []uint32 {
	out := make([]uint32, c)
	for i := range out {
		out[i] = uint32(i)
	}
	for i := range out {
		j := rand.Uint32() % c
		out[i], out[j] = out[j], out[i]
	}
	return out
}

func expectNeed(i uint32, lst []uint32) []uint32 {
	var out []uint32
	for _, idx := range lst[i+1:] {
		out = append(out, uint32(idx))
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i] < out[j]
	})
	return out
}
