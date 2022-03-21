package txtidx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoc(t *testing.T) {
	expected := "This is a test"

	di := &String{
		Len:   15,
		Start: 1,
		IndexWords: []IWordIndex{
			{
				BaseID:   0,
				Variants: []uint16{0},
				Links: []Link{
					{
						VID:  0,
						Next: MaxUint32,
					},
				},
			},
			{
				BaseID:   1,
				Variants: []uint16{1},
				Links: []Link{
					{
						VID:  0,
						Next: 2,
					},
				},
			},
		},
		UnindexWords: []UWordIndex{
			{
				ID:   1,
				Next: []uint32{3},
			},
			{
				ID:   0,
				Next: []uint32{0},
			},
		},
	}

	s := &Source{
		Indexed: []IndexedWord{
			{
				[]string{"test"},
			},
			{
				[]string{"this", "This "},
			},
		},
		Unindexed: []string{
			"a ", "is ",
		},
	}

	got := Build(s, di)
	assert.Equal(t, expected, got)

}

func TestSourceIDX(t *testing.T) {
	assert.Equal(t, SourceIDX(0), SourceIDX(1)&isUnindexed)
	assert.Equal(t, SourceIDX(1), SourceIDX(1)&rootMask)

	var r uint64 = 123
	var v uint64 = 456
	i := IndexedIDX(r, v)
	assert.True(t, i.Indexed())
	assert.Equal(t, r, i.Root())
	assert.Equal(t, v, i.Variant())
	v = 789
	assert.Equal(t, v, i.SetVariant(int(v)).Variant())

	var u int = 31415
	i = UnindexedIDX(u)
	assert.False(t, i.Indexed())
	assert.Equal(t, uint64(u), i.Unindex())
}

func TestAppendSource(t *testing.T) {
	s := NewSource()
	idx := s.UpsertIndexed("this", "This ")
	assert.Equal(t, "This ", s.Get(idx))
	assert.Equal(t, idx, s.Lookup["This "])
	assert.Equal(t, idx, s.UpsertIndexed("this", "This "))
	assert.Equal(t, idx.Root(), s.UpsertIndexed("this", "this, ").Root())

	idx = s.UpsertUnindexed("foo")
	assert.Equal(t, "foo", s.Get(idx))
	assert.Equal(t, idx, s.UpsertUnindexed("foo"))
}
