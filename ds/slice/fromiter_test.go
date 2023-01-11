package slice_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

// doneIter is an empty Iter that fulfills neither Slicer nor Lener.
type doneIter struct{}

func (doneIter) Next() (int, bool) { return 0, true }
func (doneIter) Cur() (int, bool)  { return 0, true }
func (doneIter) Done() bool        { return true }
func (doneIter) Idx() int          { return 0 }

func TestFromIterEmpty(t *testing.T) {
	assert.Nil(t, slice.FromIter[int](doneIter{}, nil))
}
