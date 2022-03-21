package txtidx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoc(t *testing.T) {
	expected := "This is a test"

	di := &DocumentIndex{
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
