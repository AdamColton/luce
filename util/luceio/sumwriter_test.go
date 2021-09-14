package luceio_test

import (
	"bytes"
	"testing"

	"github.com/adamcolton/luce/util/luceio"
	"github.com/stretchr/testify/assert"
)

func TestSumWriter(t *testing.T) {
	b, sw := luceio.BufferSumWriter()

	sw.WriteInt(123)
	sw.WriteRune('c')
	sw.WriteString("this")
	sw.WriteStrings("is", "a", "test")
	sw.Fprint("%f", 3.1415)

	n, err := sw.Rets()
	assert.NoError(t, err)
	assert.Equal(t, int64(23), n)

	assert.Equal(t, "123cthisisatest3.141500", b.String())
}

func TestSumWriterJoin(t *testing.T) {
	b, sw := luceio.BufferSumWriter()

	n, err := sw.Join([]string{"this", "is", "a", "test"}, " ")
	assert.NoError(t, err)
	assert.Equal(t, int(14), n)
	assert.Equal(t, "this is a test", b.String())
}

func TestSumWriterTo(t *testing.T) {
	out, sw := luceio.BufferSumWriter()

	s := "this is a test"
	in := bytes.NewBufferString(s)
	n, err := sw.WriterTo(in)
	assert.NoError(t, err)
	assert.Equal(t, int64(14), n)
	assert.Equal(t, s, out.String())
}
