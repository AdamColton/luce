package luceio_test

import (
	"bytes"
	"testing"

	"github.com/adamcolton/luce/util/luceio"
	"github.com/stretchr/testify/assert"
)

func TestStringWriterTo(t *testing.T) {
	expected := "testing"
	str := luceio.StringWriterTo(expected)
	buf := bytes.NewBuffer(nil)
	str.WriteTo(buf)

	assert.Equal(t, expected, buf.String())
}
