// Package sliceidx provides SliceIdx as a building block for creating
// types that fulfill idx.Index.
package sliceidx

// SliceIdx holds the logic for the slice index.
type SliceIdx struct {
	SliceLen int
	MaxIdx   int
	Recycled []int
}

// New creates a SliceIdx and sets the sliceLen.
func New(sliceLen int) SliceIdx {
	return SliceIdx{
		SliceLen: sliceLen,
	}
}

// NextIdx gets the next index for the slice and if it requires an append
// action.
func (si *SliceIdx) NextIdx() (idx int, app bool) {
	if ln := len(si.Recycled); ln > 0 {
		idx = si.Recycled[ln-1]
		si.Recycled = si.Recycled[:ln-1]
	} else {
		idx = si.MaxIdx
		si.MaxIdx++
		app = si.MaxIdx > si.SliceLen
		if app {
			si.SliceLen = si.MaxIdx
		}
	}
	return
}

// SetSliceLen will set the SliceLen to newlen if newlen > SliceLen.
func (si *SliceIdx) SetSliceLen(newlen int) {
	if newlen > si.SliceLen {
		si.SliceLen = newlen
	}
}

// Recycle appends the idx to Recycled.
func (si *SliceIdx) Recycle(idx int) {
	si.Recycled = append(si.Recycled, idx)
}
