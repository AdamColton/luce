package gothicgo

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/testify/assert"
)

type writerPrepper struct {
	str      string
	prepErr  error
	writeErr error
	closed   bool
}

func (wp *writerPrepper) WriteTo(w io.Writer) (int64, error) {
	b := []byte(wp.str)
	w.Write(b)
	return int64(len(b)), wp.writeErr
}

func (wp *writerPrepper) Prepare() error {
	return wp.prepErr
}

type wrapBuffer struct {
	buf      *bytes.Buffer
	closed   bool
	closeErr error
	writeErr error
}

func (mc *wrapBuffer) Write(b []byte) (int, error) {
	i, err := mc.buf.Write(b)
	if err != nil {
		return i, err
	}
	return i, mc.writeErr
}

func (mc *wrapBuffer) Close() error {
	mc.closed = true
	return mc.closeErr
}

func TestFile(t *testing.T) {
	ctx := NewMemoryContext()

	// Use wrapBuffer to confirm that File.Generate checks if the Writer is a
	// Closer and calls close.
	cf := ctx.CreateFile
	var wb *wrapBuffer
	ctx.CreateFile = func(path string) (io.Writer, error) {
		w, err := cf(path)
		wb = &wrapBuffer{buf: w.(*bytes.Buffer)}
		return wb, err
	}
	defer func() {
		assert.True(t, wb.closed)
	}()

	pkg := ctx.MustPackage("foo")
	file := pkg.File("bar")
	wp := &writerPrepper{str: `var test = "testing file generation"`}
	file.AddWriterTo(wp)

	assert.Equal(t, pkg, file.Package())
	assert.Equal(t, "bar", file.Name())

	assert.NoError(t, ctx.Export())

	str := ctx.Last.String()
	assert.Contains(t, str, `var test = "testing file generation"`)
	assert.Contains(t, str, "package foo")

}

func TestFilePrepErr(t *testing.T) {
	ctx := NewMemoryContext()
	pkg := ctx.MustPackage("foo")
	file := pkg.File("bar")
	file.AddWriterTo(&writerPrepper{prepErr: fmt.Errorf("testing file prep error")})
	err := ctx.Export()

	assert.Equal(t, "Prepare package foo: While preparing file bar: testing file prep error", err.Error())
}

func TestFileWriteToErr(t *testing.T) {
	ctx := NewMemoryContext()
	pkg := ctx.MustPackage("foo")
	file := pkg.File("bar")
	file.AddWriterTo(&writerPrepper{writeErr: fmt.Errorf("testing file write error")})
	err := ctx.Export()

	assert.Equal(t, "Generate package foo: WriteTo Error in Generate file foo/bar: testing file write error", err.Error())
}

func TestFileCloseErr(t *testing.T) {
	ctx := NewMemoryContext()

	// Use wrapBuffer to cause error closing the writer
	cf := ctx.CreateFile
	ctx.CreateFile = func(path string) (io.Writer, error) {
		w, err := cf(path)
		return &wrapBuffer{
			buf:      w.(*bytes.Buffer),
			closeErr: fmt.Errorf("Close Error"),
		}, err
	}

	pkg := ctx.MustPackage("foo")
	file := pkg.File("bar")
	wp := &writerPrepper{str: `var test = "testing file close error"`}
	file.AddWriterTo(wp)

	err := ctx.Export()
	assert.Equal(t, "Generate package foo: Closing writer for file foo/bar: Close Error", err.Error())
}

func TestFileFormatErr(t *testing.T) {
	ctx := NewMemoryContext()
	pkg := ctx.MustPackage("foo")
	file := pkg.File("bar")
	file.AddWriterTo(&writerPrepper{str: "testing file format error"})
	err := ctx.Export()

	assert.Equal(t, "Generate package foo: Failed to format foo/bar:: 4:1: expected declaration, found testing", err.Error())
	assert.Contains(t, ctx.Last.String(), "testing file format error")
}

func TestFileDoubleGet(t *testing.T) {
	ctx := NewMemoryContext()
	pkg := ctx.MustPackage("foo")
	file := pkg.File("bar")
	again := pkg.File("bar")

	assert.Equal(t, file, again)
}
