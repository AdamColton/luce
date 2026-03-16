package morph_test

import (
	"strconv"
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/morph"
	"github.com/adamcolton/luce/lerr"
	"github.com/stretchr/testify/assert"
)

func TestMorphSlice(t *testing.T) {
	m := lmap.New(map[string]int{
		"1": 1,
		"2": 2,
		"3": 3,
		"4": 4,
		"5": 5,
	})
	e := morph.NewKVToKV(func(k string, v int) (uint64, float64) {
		return lerr.Must(strconv.ParseUint(k, 10, 64)), float64(v)
	}).Eacher(m)
	out := lmap.FromEacher(e)
	expected := lmap.New(map[uint64]float64{
		1: 1,
		2: 2,
		3: 3,
		4: 4,
		5: 5,
	})
	assert.Equal(t, expected, out)
}
