package lmap_test

import (
	"strconv"
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/math/numiter"
	"github.com/stretchr/testify/assert"
)

func TestIter(t *testing.T) {
	i := numiter.NewRange(0, 10, 1).Wrap().Iter()
	m := lmap.FromIter(i, func(i, idx int) (int, string, bool) {
		return i, strconv.Itoa(i), true
	})
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

	empty := lmap.FromIter(i, func(i, idx int) (int, string, bool) {
		return 0, "", false
	})
	assert.Equal(t, 0, empty.Len())
}
