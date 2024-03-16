package numiter_test

import (
	"testing"

	"github.com/adamcolton/luce/math/numiter"
	"github.com/stretchr/testify/assert"
)

func TestWrapIter(t *testing.T) {
	r := numiter.NewRange(0, 3, 1)
	assert.Equal(t, []int{0, 1, 2}, []int(r.Wrap().Slice(nil)))

	var got []int
	it := r.Iter()
	for v, done := it.Cur(); !done; v, done = it.Next() {
		got = append(got, v)
	}
	assert.Equal(t, []int{0, 1, 2}, got)
}
