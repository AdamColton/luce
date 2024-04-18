package lmap_test

import (
	"strconv"
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/morph"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

func TestIter(t *testing.T) {
	i := slice.New([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}).Iter()
	mi := morph.NewValAll(func(i int) lmap.KeyVal[int, string] {
		return lmap.NewKV(i, strconv.Itoa(i))
	}).Iter(i)
	m := lmap.FromIter(mi)
	expected := lmap.New(map[int]string{
		0: "0",
		1: "1",
		2: "2",
		3: "3",
		4: "4",
		5: "5",
		6: "6",
		7: "7",
		8: "8",
		9: "9",
	})
	assert.Equal(t, expected, m)

	mi = morph.NewVal(func(i int) (lmap.KeyVal[int, string], bool) {
		return lmap.NewKV(0, ""), false
	}).Iter(i)
	empty := lmap.FromIter(mi)
	assert.Nil(t, empty.Mapper)
}
