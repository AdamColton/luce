package luceio_test

import (
	"testing"

	"github.com/adamcolton/luce/util/luceio"
	"github.com/stretchr/testify/assert"
)

func TestSumWriterJoinEmpty(t *testing.T) {
	buf, sw := luceio.BufferSumWriter()

	n, err := sw.Join(nil, ", ")
	assert.NoError(t, err)
	assert.Zero(t, n)

	n, err = sw.Join([]string{}, ", ")
	assert.NoError(t, err)
	assert.Zero(t, n)
	assert.Zero(t, buf.Len())
}
