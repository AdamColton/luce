package numiter_test

import (
	"testing"

	"github.com/adamcolton/luce/math/cmpr/cmprtest"
	"github.com/adamcolton/luce/math/numiter"
)

func TestNGrid(t *testing.T) {
	rng := numiter.NewRange(0, 2, 1)
	expected := [][]int{{0, 0}, {1, 0}, {0, 1}, {1, 1}}
	cmprtest.Equal(t, numiter.NGrid(rng, 2), expected)
}
