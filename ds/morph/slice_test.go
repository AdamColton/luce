package morph_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/morph"
	"github.com/stretchr/testify/assert"
)

func TestValSlice(t *testing.T) {
	in := []int{1, 2, 3, 4, 5}
	assert.Equal(t, []string{"n2", "n4"}, []string(even.Slice(in, nil)))
	assert.Nil(t, even.Slice([]int{1, 3}, nil))
	assert.Nil(t, even.Slice(nil, nil))

	// The buffer is used when it has the capacity.
	buf := make([]string, 0, 4)
	got := even.Slice(in, buf)
	assert.Same(t, &buf[:1][0], &got[0])

	vals, _ := collect(even.SliceIter(in))
	assert.Equal(t, []string{"n2", "n4"}, vals)

	assert.Equal(t, []string{"n2", "n4"}, []string(even.IterSlice(newIter(1, 2, 4), nil)))
}

func TestValAllSlice(t *testing.T) {
	double := morph.NewValAll(func(i int) int { return i * 2 })
	in := []int{1, 2, 3}
	assert.Equal(t, []int{2, 4, 6}, []int(double.Slice(in, nil)))

	buf := make([]int, 0, 4)
	got := double.Slice(in, buf)
	assert.Same(t, &buf[:1][0], &got[0])

	it := double.SliceIter(in)
	var vals []int
	for v, done := it.Cur(); !done; v, done = it.Next() {
		vals = append(vals, v)
	}
	assert.Equal(t, []int{2, 4, 6}, vals)

	assert.Equal(t, []int{2, 4}, []int(double.IterSlice(newIter(1, 2), nil)))
}
