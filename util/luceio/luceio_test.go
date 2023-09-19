package luceio_test

import (
	"bytes"
	"testing"

	"github.com/adamcolton/luce/util/luceio"
	"github.com/stretchr/testify/assert"
)

func TestStringWriter(t *testing.T) {
	buf := &bytes.Buffer{}
	var sw luceio.StringWriter = buf
	_, err := sw.WriteString("testing")
	assert.NoError(t, err)
	assert.Equal(t, "testing", buf.String())
}

func TestStringsWriter(t *testing.T) {
	buf := &bytes.Buffer{}
	var sw luceio.StringsWriter = luceio.NewSumWriter(buf)
	n, err := sw.WriteStrings("this", "is", "a", "test")
	assert.NoError(t, err)
	assert.Equal(t, 11, n)
	assert.Equal(t, "thisisatest", buf.String())
}
