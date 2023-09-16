package sliceidx_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/idx/sliceidx"
	"github.com/stretchr/testify/assert"
)

func TestSliceIdx(t *testing.T) {
	si := sliceidx.New(10)
	assert.Equal(t, 10, si.SliceLen)
	for i := 0; i < 15; i++ {
		idx, app := si.NextIdx()
		assert.Equal(t, i, idx)
		assert.Equal(t, i >= 10, app)
	}
	si.Recycle(5)
	idx, app := si.NextIdx()
	assert.Equal(t, 5, idx)
	assert.False(t, app)

	si.SetSliceLen(16)
	idx, app = si.NextIdx()
	assert.Equal(t, 15, idx)
	assert.False(t, app)

	idx, app = si.NextIdx()
	assert.Equal(t, 16, idx)
	assert.True(t, app)
}
