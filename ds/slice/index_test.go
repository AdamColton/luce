package slice_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	s := slice.Slice[int]{3, 1, 4, 1, 5, 9, 2, 6, 5, 3}
	i := slice.NewIndex(3, 2)
	assert.Equal(t, slice.Index{3, 5}, i)
	assert.Equal(t, s[3:5], s.Sub(i))
	assert.Equal(t, 4, i.Last())
	assert.Equal(t, 2, i.Len())
	assert.Equal(t, 4, i.AtIdx(1))
	i = i.Next(4)
	assert.Equal(t, slice.Index{5, 9}, i)
	assert.Equal(t, 6, i.AtIdx(1))

	bs := slice.IdxMake[bool](i)
	assert.Len(t, bs, i[1])
}
