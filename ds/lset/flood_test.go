package lset_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/ds/lset"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/math/numiter"
	"github.com/adamcolton/luce/util/filter"
	"github.com/stretchr/testify/assert"
)

func TestFlood(t *testing.T) {

	type Pt [2]int
	distFn := func(pt Pt) int { return pt[0]*pt[0] + pt[1]*pt[1] }
	f := filter.New(func(pt Pt) bool {
		return distFn(pt) < 100
	})
	less := func(i, j Pt) bool {
		return i[0] < j[0] || (i[0] == j[0] && i[1] <= j[1])
	}

	dirs := slice.Slice[Pt]([]Pt{{1, 0}, {0, 1}, {-1, 0}, {0, -1}})
	circleProc := func(pt Pt, add func(Pt)) {
		for _, d := range dirs {
			dpt := [2]int{pt[0] + d[0], pt[1] + d[1]}
			if f(dpt) {
				//fmt.Println(dpt)
				add(dpt)
			}
		}
	}

	tfrm := list.TransformAny(func(in []int) Pt {
		return Pt{in[0], in[1]}
	})
	i := f.Iter(tfrm.New(numiter.Grid(-10, 10, 1, -10, 10, 1)).Iter())
	expected := slice.FromIter(i, nil)
	expected.Sort(less)

	got := lset.Flood(circleProc, Pt{0, 0}, Pt{2, 2}).Slice(nil).Sort(less)

	assert.Equal(t, expected, got)
}
