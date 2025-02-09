package lmap_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/stretchr/testify/assert"
)

func TestTransform(t *testing.T) {
	tf := lmap.NewTransformFunc(func(k int, s string) (string, int, bool) {
		return s, k, s != ""
	})

	m := map[int]string{
		1:  "1",
		2:  "2",
		3:  "",
		4:  "4",
		5:  "5",
		6:  "6",
		7:  "",
		8:  "8",
		9:  "9",
		10: "10",
	}

	expected := map[string]int{
		"1":  1,
		"2":  2,
		"4":  4,
		"5":  5,
		"6":  6,
		"8":  8,
		"9":  9,
		"10": 10,
	}

	got := lmap.TransformMap(m, tf).Map()
	assert.Equal(t, expected, got)

	buf := lmap.New(map[string]int{
		"100": 100,
	})
	expected["100"] = 100
	got = lmap.Transform(lmap.New(m), buf, tf).Map()
	assert.Equal(t, expected, got)
}
