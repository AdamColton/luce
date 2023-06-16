package huffslice

import (
	"testing"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

func TestHuffSlice(t *testing.T) {
	tt := map[string][]uint32{
		"basic": {
			3, 1, 4, 5, 9, 2, 6, 5, 3, 5, 8, 9, 7, 9, 3, 2, 3, 8, 4, 6, 2, 6, 4,
			3, 3, 8, 3, 2, 7, 9, 5, 0},
		"allUnique":     {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		"withTokenOnce": {1, 1, 1, 2, 2, 2, 1000, 3, 4},
		"withToken":     {1, 1, 1, 2, 2, 2, 1000, 1000, 3, 4},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			ue := NewEncoder(len(tc), uint32(1000))
			ue.Slice = append(ue.Slice, tc...)

			u := ue.Encode()
			got := slice.IterSlice(u.Iter(), nil)

			assert.Equal(t, tc, got)
		})
	}
}
