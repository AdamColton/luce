package slice_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

func TestClone(t *testing.T) {
	data := []int{3, 1, 4, 1, 5, 9}
	cp := slice.Clone(data)
	assert.Equal(t, data, cp)
	data[0] = 0
	assert.Equal(t, 3, cp[0])
}

func TestSwap(t *testing.T) {
	data := []int{3, 1, 4, 1, 5, 9}
	slice.Swap(data, 0, 1)
	assert.Equal(t, 1, data[0])
	assert.Equal(t, 3, data[1])

}

func TestKeys(t *testing.T) {
	data := map[int]string{
		1: "1",
		2: "2",
		3: "3",
		4: "4",
		5: "5",
		6: "6",
	}
	got := slice.Keys(data)
	slice.Less[int](func(i, j int) bool {
		return i < j
	}).Sort(got)
	expected := []int{1, 2, 3, 4, 5, 6}
	assert.Equal(t, expected, got)
}

func TestVals(t *testing.T) {
	data := map[int]string{
		1: "1",
		2: "2",
		3: "3",
		4: "4",
		5: "5",
		6: "6",
	}
	got := slice.Vals(data)
	slice.Less[string](func(i, j string) bool {
		return i < j
	}).Sort(got)
	expected := []string{"1", "2", "3", "4", "5", "6"}
	assert.Equal(t, expected, got)
}

func TestLess(t *testing.T) {
	i := []int{6, 7, 9, 2, 3, 4, 1, 5, 8}
	slice.Less[int](func(i, j int) bool {
		return i < j
	}).Sort(i)
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	assert.Equal(t, expected, i)
}
