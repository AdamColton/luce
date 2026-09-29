package morph_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/morph"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

func TestDeref(t *testing.T) {
	i := 3
	assert.Equal(t, 3, morph.Deref(&i))

	assert.Panics(t, func() {
		morph.Deref[int](nil)
	})
}

func TestNewDeref(t *testing.T) {
	a, b, c := 3, 1, 4
	ptrs := []*int{&a, &b, &c}
	got := morph.NewDeref[int]().Slice(ptrs, nil)
	assert.Equal(t, slice.Slice[int]{3, 1, 4}, got)
}
