package cmprtest_test

import (
	"testing"

	"github.com/adamcolton/luce/math/cmpr/cmprtest"
)

func TestCmprTest(t *testing.T) {
	d := float64(cmprtest.Small / 10)
	a := 3.1415
	b := a + d
	cmprtest.Equal(t, a, b)
}
