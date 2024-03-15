package lerr_test

import (
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/stretchr/testify/assert"
)

func TestSliceErrs(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	b := []int{1, 200, 3, 4}

	fn := func(i int) error {
		ai, bi := a[i], b[i]
		return lerr.NewNotEqual(ai == bi, ai, bi)
	}

	err := lerr.NewSliceErrs(len(a), len(b), fn)
	assert.Equal(t, "Lengths do not match: Expected 5 got 4\n\t1: Expected 2 got 200", err.Error())

	b[1] = 2
	b = append(b, 5)
	err = lerr.NewSliceErrs(len(a), len(b), fn)
	assert.NoError(t, err)

	var idxs []int
	fn = func(i int) error {
		idxs = append(idxs, i)
		return nil
	}
	err = lerr.NewSliceErrs(3, -1, fn)
	assert.NoError(t, err)
	assert.Equal(t, []int{0, 1, 2}, idxs)

	var s lerr.SliceErrs
	s = s.AppendF(3, "test %d", 123)
	s = s.Append(10, lerr.Str("foo test"))
	s = s.Append(10, lerr.Str("bar test"))
	s = s.Append(10, lerr.Str("baz test"))

	restore := lerr.MaxSliceErrs
	lerr.MaxSliceErrs = 2
	assert.Equal(t, "\t3: test 123\n\t10: foo test\nOmitting 2 more", s.Error())
	lerr.MaxSliceErrs = restore
}
