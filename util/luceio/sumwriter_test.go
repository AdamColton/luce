package luceio_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
	"github.com/adamcolton/luce/util/upgrade"
	"github.com/stretchr/testify/assert"
)

type errWriter struct {
	after int
	err   error
}

func (ew *errWriter) Write(b []byte) (int, error) {
	if ew.after <= 0 {
		return 0, ew.err
	}
	ln := len(b)
	ew.after -= ln
	return ln, nil
}

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

func TestWriteStringsErr(t *testing.T) {
	ew := &errWriter{
		after: 10,
		err:   lerr.Str("test err"),
	}
	sw := luceio.NewSumWriter(ew)
	n, err := sw.WriteStrings("this", "is", "another", "test")
	assert.True(t, n >= 10)
	assert.Equal(t, ew.err, err)
}

func TestWriteErr(t *testing.T) {
	ew := &errWriter{
		after: 10,
		err:   lerr.Str("test err"),
	}
	sw := luceio.NewSumWriter(ew)
	n, err := sw.WriteString("0123456789")
	assert.Equal(t, 10, n)
	assert.NoError(t, err)

	n, err = sw.WriteString("err returned")
	assert.Equal(t, 0, n)
	assert.Equal(t, ew.err, err)

	sw.Err = nil
	sw.Cache = []byte("cache")
	n, err = sw.WriteString("cache err")
	assert.Equal(t, 0, n)
	assert.Equal(t, ew.err, err)

	n, err = sw.WriteString("init err check")
	assert.Equal(t, 0, n)
	assert.Equal(t, ew.err, err)
}

func TestSumWriterJoin(t *testing.T) {
	b, sw := luceio.BufferSumWriter()

	n, err := sw.Join([]string{"this", "is", "a", "test"}, " ")
	assert.NoError(t, err)
	assert.Equal(t, int(14), n)
	assert.Equal(t, "this is a test", b.String())
}

// strIter is a minimal liter.Iter[string] over a slice.
type strIter struct {
	strs []string
	idx  int
}

func (si *strIter) Cur() (string, bool) {
	if si.idx >= len(si.strs) {
		return "", true
	}
	return si.strs[si.idx], false
}

func (si *strIter) Next() (string, bool) {
	si.idx++
	return si.Cur()
}

func (si *strIter) Done() bool {
	return si.idx >= len(si.strs)
}

func (si *strIter) Idx() int {
	return si.idx
}

func TestSumWriterIter(t *testing.T) {
	b, sw := luceio.BufferSumWriter()
	n, err := sw.Iter(&strIter{strs: []string{"this", "is", "a", "test"}}, " ")
	assert.NoError(t, err)
	assert.Equal(t, 14, n)
	assert.Equal(t, "this is a test", b.String())

	// n only counts bytes written by this call
	n, err = sw.Iter(&strIter{strs: []string{"", "x", "y"}}, "-")
	assert.NoError(t, err)
	assert.Equal(t, 4, n)
	assert.Equal(t, "this is a test-x-y", b.String())
	assert.Equal(t, int64(18), sw.Sum)

	// empty iterator writes nothing
	b, sw = luceio.BufferSumWriter()
	n, err = sw.Iter(&strIter{}, " ")
	assert.NoError(t, err)
	assert.Equal(t, 0, n)
	assert.Equal(t, "", b.String())

	// starts from the current value, not reset
	b, sw = luceio.BufferSumWriter()
	it := &strIter{strs: []string{"a", "b", "c"}, idx: 1}
	n, err = sw.Iter(it, ",")
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, "b,c", b.String())

	// errors are returned
	expected := lerr.Str("test err")
	sw = luceio.NewSumWriter(&errWriter{after: 5, err: expected})
	_, err = sw.Iter(&strIter{strs: []string{"this", "is", "a", "test"}}, " ")
	assert.Equal(t, expected, err)
}

func TestWriterToErr(t *testing.T) {
	sw := luceio.NewSumWriter(bytes.NewBuffer(nil))
	sw.Err = lerr.Str("test err")
	n, err := sw.WriterTo(bytes.NewBufferString("this is a test"))
	assert.Equal(t, int64(0), n)
	assert.Equal(t, sw.Err, err)
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

func TestSumWriterCache(t *testing.T) {
	out, sw := luceio.BufferSumWriter()
	sw.WriteString("this")
	sw.AppendCacheString("\n")
	sw.WriteString("is a test")
	assert.NoError(t, sw.Err)
	assert.Equal(t, int64(14), sw.Sum)
	assert.Equal(t, "this\nis a test", out.String())
}

func TestSumWriterUpgrade(t *testing.T) {
	sw := luceio.NewSumWriter(bytes.NewBufferString("stringer"))
	str, is := upgrade.To[fmt.Stringer](sw)
	assert.True(t, is)
	assert.Equal(t, "stringer", str.String())
}
