package morph_test

import (
	"strconv"
	"testing"

	"github.com/adamcolton/luce/ds/morph"
	"github.com/stretchr/testify/assert"
)

// ints is a list.List of ints.
type ints []int

func (i ints) AtIdx(idx int) int { return i[idx] }
func (i ints) Len() int          { return len(i) }

func TestList(t *testing.T) {
	calls := 0
	va := morph.NewValAll(func(i int) string {
		calls++
		return "n" + strconv.Itoa(i)
	})
	l := va.List(ints{1, 2, 3})
	assert.Equal(t, 3, l.Len())
	assert.Equal(t, "n1", l.AtIdx(0))
	assert.Equal(t, "n3", l.AtIdx(2))

	// Nothing is cached, so every AtIdx calls the ValAll.
	assert.Equal(t, 2, calls)
	l.AtIdx(0)
	assert.Equal(t, 3, calls)

	direct := morph.List[int, string]{List: ints{4, 5}, ValAll: va}
	assert.Equal(t, "n5", direct.AtIdx(1))
}
