package cmpr_test

import (
	"testing"

	"github.com/adamcolton/luce/math/cmpr"
	"github.com/stretchr/testify/assert"
)

func TestMin(t *testing.T) {
	assert.Equal(t, 5, cmpr.Min(5, 6))
	assert.Equal(t, 6, cmpr.Min(7, 6))
	assert.Equal(t, "good-bye", cmpr.Min("hello", "good-bye"))
}

func TestMinN(t *testing.T) {
	assert.Equal(t, 1, cmpr.MinN(3, 1, 4, 1, 5, 9))
}

func TestMax(t *testing.T) {
	assert.Equal(t, 6, cmpr.Max(5, 6))
	assert.Equal(t, 7, cmpr.Max(7, 6))
	assert.Equal(t, "hello", cmpr.Max("hello", "good-bye"))
}

func TestMaxN(t *testing.T) {
	assert.Equal(t, 9, cmpr.MaxN(3, 1, 4, 1, 5, 9))
	assert.Equal(t, 0, cmpr.MaxN[int]())
}
