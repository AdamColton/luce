// Package sliceidx provides SliceIdx as a building block for creating
// types that fulfill idx.Index.
package sliceidx

// SliceIdx holds the logic for the slice index. It hands out positions in a
// slice, reusing positions that were recycled before growing.
type SliceIdx struct {
	// SliceLen is the length the slice needs to hold every index handed out.
	SliceLen int
	// MaxIdx is the next index that has never been handed out.
	MaxIdx int
	// Recycled holds indexes that are free to be handed out again.
	Recycled []int
}

// New creates a SliceIdx and sets the sliceLen.
func New(sliceLen int) SliceIdx {
	return SliceIdx{
		SliceLen: sliceLen,
	}
}

// NextIdx gets the next index for the slice and if it requires an append
// action. The most recently recycled index is used first and does not require
// an append. Otherwise it returns the next unused index, and app is true if
// that index is past SliceLen, in which case SliceLen is raised to include it.
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

// Recycle appends the idx to Recycled so NextIdx will hand it out again. It
// does not check that idx is in use or has not already been recycled.
func (si *SliceIdx) Recycle(idx int) {
	si.Recycled = append(si.Recycled, idx)
}
