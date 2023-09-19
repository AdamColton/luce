package luceio_test

import (
	"bytes"
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
	"github.com/stretchr/testify/assert"
)

func TestWriterTos(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	tos := luceio.WriterTos{luceio.StringWriterTo("test1"), luceio.StringWriterTo("test2")}
	tos = luceio.WriterTos{tos, luceio.StringWriterTo("test3")}

	assert.Len(t, tos, 2)
	tos = tos.Merge()
	assert.Len(t, tos, 3)

	n, err := tos.WriteTo(buf)
	assert.Equal(t, "test1test2test3", buf.String())
	assert.NoError(t, err)
	assert.Equal(t, int64(15), n)

	buf.Reset()
	n, err = tos.Seperator(":").WriteTo(buf)
	assert.NoError(t, err)
	assert.Equal(t, int64(17), n)
	assert.Equal(t, "test1:test2:test3", buf.String())

	buf.Reset()
	n, err = tos.Seperator([]byte{'|'}).WriteTo(buf)
	assert.NoError(t, err)
	assert.Equal(t, int64(17), n)
	assert.Equal(t, "test1|test2|test3", buf.String())
}

func TestWriterTosErr(t *testing.T) {
	ew := &errWriter{
		after: 10,
		err:   lerr.Str("test err"),
	}

	tos := luceio.WriterTos{
		luceio.StringWriterTo("test1"),
		luceio.StringWriterTo("test2"),
		luceio.StringWriterTo("test3"),
	}
	n, err := tos.WriteTo(ew)
	assert.Equal(t, int64(10), n)
	assert.Equal(t, ew.err, err)

	// trigger error check after writing from WriterTo
	ew.after = 12
	n, err = tos.Seperator(":").WriteTo(ew)
	assert.Equal(t, int64(12), n)
	assert.Equal(t, ew.err, err)

	// trigger error check after writing seperator
	ew.after = 11
	n, err = tos.Seperator(":").WriteTo(ew)
	assert.Equal(t, int64(11), n)
	assert.Equal(t, ew.err, err)
}
