package merkle

import (
	"crypto/sha256"
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

		done, t2 := a.Done()
		assert.False(t, done)
		for i, idx := range lst {
			done, t2 = a.Done()
			assert.False(t, done)
			l = t1.Leaf(idx)
			assert.True(t, a.Add(l))
			assert.Equal(t, td, l.Digest(sha256.New()))

			assert.Equal(t, expectNeed(i, lst), a.Need())

			if i == 0 {
				// test incomplete tree calls
				incpTr := Tree(a.root)
				assert.Equal(t, -1, incpTr.size())
				assert.Equal(t, -1, incpTr.Count())
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
	}
}

func randList(c int) []int {
	out := make([]int, c)
	for i := range out {
		out[i] = i
	}
	for i := range out {
		j := rand.Intn(c)
		out[i], out[j] = out[j], out[i]
	}
	return out
}

func expectNeed(i int, lst []int) []uint32 {
	var out []uint32
	for _, idx := range lst[i+1:] {
		out = append(out, uint32(idx))
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i] < out[j]
	})
	return out
}
