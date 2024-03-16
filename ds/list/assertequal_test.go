package list_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/lerr"
	"github.com/stretchr/testify/assert"
)

func TestAssertEqualWrongType(t *testing.T) {
	w := list.Slice([]float64{3, 1, 4})
	to := "not a list"
	err := w.AssertEqual(to, 1e-6)
	assert.Equal(t, lerr.NewTypeMismatch(w, to), err)

	// A list of another type is not a List[float64].
	err = w.AssertEqual(list.Slice([]int{3, 1, 4}), 1e-6)
	assert.Error(t, err)
}
