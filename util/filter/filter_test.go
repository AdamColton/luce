package filter_test

import (
	"testing"

	"github.com/adamcolton/luce/util/filter"
	"github.com/stretchr/testify/assert"
)

func TestFilterOps(t *testing.T) {
	gt4 := filter.Filter[int](func(i int) bool {
		return i > 4
	})
	lt8 := filter.Filter[int](func(i int) bool {
		return i < 8
	})

	gt4AndLt8 := gt4.And(lt8)
	assert.False(t, gt4AndLt8(3))
	assert.True(t, gt4AndLt8(5))
	assert.False(t, gt4AndLt8(10))

	lte4OrGte8 := gt4.Not().Or(lt8.Not())
	assert.True(t, lte4OrGte8(3))
	assert.False(t, lte4OrGte8(5))
	assert.True(t, lte4OrGte8(10))
}
