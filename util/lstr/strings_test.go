package lstr_test

import (
	"strings"
	"testing"
	"time"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/iter"
	"github.com/adamcolton/luce/util/lstr"
	"github.com/stretchr/testify/assert"
)

func TestStringsIter(t *testing.T) {
	s := lstr.NewStrings([]string{"", "this ", "is ", " ", " a ", "test"})
	var i iter.Iter[string] = s
	expect := slice.NewIter([]string{"this", "is", "a", "test"})
	expectIdx := slice.NewIter([]int{1, 2, 4, 5})
	iter.ForIdx(i, func(str string, sIdx int) {
		assert.Equal(t, expect.Pop(), str)
		assert.Equal(t, expectIdx.Pop(), sIdx)
	})

	i.(iter.Starter[string]).Start()
	str, _ := i.Cur()
	assert.Equal(t, "this", str)
	assert.Equal(t, 1, i.Idx())
}

func TestStringsSub(t *testing.T) {
	s := lstr.NewStrings(strings.Split("1,2,3,4|,,|,5,6,7|8,9,10", "|"))
	expect := slice.NewIter([][]string{
		{"1", "2", "3", "4"},
		{"5", "6", "7"},
		{"8", "9", "10"},
	})
	for !s.Done() {
		e := slice.NewIter(expect.Pop())
		iter.For(s.Sub(","), func(str string) {
			assert.Equal(t, e.Pop(), str)
		})
	}
	assert.Nil(t, s.Sub(","))
}

func TestStringsFloat64(t *testing.T) {
	s := lstr.NewStrings(strings.Split(",3.1415,1.414,,1.618,", ","))
	expect := slice.NewIter([]float64{3.1415, 1.414, 1.618})
	for !s.Done() {
		assert.Equal(t, expect.Pop(), s.Float64())
	}
	s.Err = lerr.Str("test error")
	assert.Equal(t, 0.0, s.Float64())
}

func TestStringsInt(t *testing.T) {
	s := lstr.NewStrings(strings.Split(",,3,1,4,,1,5,9,2,6,,,", ","))
	expect := slice.NewIter([]int{3, 1, 4, 1, 5, 9, 2, 6})
	for !s.Done() {
		assert.Equal(t, expect.Pop(), s.Int())
	}
	s.Err = lerr.Str("test error")
	assert.Equal(t, 0, s.Int())
}

func TestStringsDate(t *testing.T) {
	layout := "2006_01_02"
	s := lstr.NewStrings(strings.Split("1984_07_12 , 2017_07_03", ","))
	expect := slice.NewIter([]time.Time{
		lerr.Must(time.Parse(layout, "1984_07_12")),
		lerr.Must(time.Parse(layout, "2017_07_03")),
	})
	for !s.Done() {
		assert.Equal(t, expect.Pop(), s.Date(layout))
	}
	s.Err = lerr.Str("test error")
	var defaultTime time.Time
	assert.Equal(t, defaultTime, s.Date(layout))
}
