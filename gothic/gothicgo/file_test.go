package gothicgo

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/testify/assert"
)

type writerPrepperNamer struct {
	str      string
	prepErr  error
	writeErr error
	closed   bool
	name     string
}

func (wp *writerPrepperNamer) WriteTo(w io.Writer) (int64, error) {
	b := []byte(wp.str)
	w.Write(b)
	return int64(len(b)), wp.writeErr
}

func (wp *writerPrepperNamer) Prepare() error {
	return wp.prepErr
}

func (wp *writerPrepperNamer) ScopeName() string {
	return wp.name
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
	wp := &writerPrepperNamer{
		str:  `var test = "testing file generation"`,
		name: "test",
	}
	assert.NoError(t, file.AddWriterTo(wp))

	assert.Equal(t, pkg, file.Package())
	assert.Equal(t, "bar", file.Name())

	ctx.MustExport()

	assert.Contains(t, ctx.Last(), `var test = "testing file generation"`)
	assert.Contains(t, ctx.Last(), "package foo")

}

func TestFilePrepErr(t *testing.T) {
	ctx, file := newFile("foo", "bar")
	assert.NoError(t, file.AddWriterTo(&writerPrepperNamer{prepErr: fmt.Errorf("testing file prep error")}))
	err := ctx.Export()

	assert.Equal(t, "Prepare package foo: While preparing file bar: testing file prep error", err.Error())
}

func TestFileWriteToErr(t *testing.T) {
	ctx, file := newFile("foo", "bar")
	assert.NoError(t, file.AddWriterTo(&writerPrepperNamer{writeErr: fmt.Errorf("testing file write error")}))
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
	wp := &writerPrepperNamer{str: `var test = "testing file close error"`}
	assert.NoError(t, file.AddWriterTo(wp))

	err := ctx.Export()
	assert.Equal(t, "Generate package foo: Closing writer for file foo/bar: Close Error", err.Error())
}

func TestFileFormatErr(t *testing.T) {
	ctx, file := newFile("foo", "bar")
	assert.NoError(t, file.AddWriterTo(&writerPrepperNamer{str: "testing file format error"}))
	err := ctx.Export()

	assert.Equal(t, "Generate package foo: Failed to format foo/bar:: 4:1: expected declaration, found testing", err.Error())
	assert.Contains(t, ctx.Last(), "testing file format error")
}

func TestFileDoubleGet(t *testing.T) {
	ctx := NewMemoryContext()
	pkg := ctx.MustPackage("foo")
	file := pkg.File("bar")
	again := pkg.File("bar")

	assert.Equal(t, file, again)
}

func TestNamerCollision(t *testing.T) {
	ctx := NewMemoryContext()
	pkg := ctx.MustPackage("foo")

	bar := pkg.File("bar")
	barWP := &writerPrepperNamer{
		str:  `var test = "testing name collision 1"`,
		name: "test",
	}
	assert.NoError(t, bar.AddWriterTo(barWP))

	baz := pkg.File("baz")
	bazWP := &writerPrepperNamer{
		str:  `var test = "testing name collision 2"`,
		name: "test",
	}
	err := baz.AddWriterTo(bazWP)
	assert.Equal(t, "File.AddWriterTo: Name 'test' already exists in package 'foo'", err.Error())
}

func TestNamerRename(t *testing.T) {
	ctx := NewMemoryContext()
	pkg := ctx.MustPackage("foo")

	bar := pkg.File("bar")
	barWP := &writerPrepperNamer{
		str:  `var test = "testing update namer 1"`,
		name: "test",
	}
	assert.NoError(t, bar.AddWriterTo(barWP))
	barWP.name = "rename"
	barWP.str = `var rename = "testing update namer 1.1"`
	assert.NoError(t, bar.UpdateNamer(barWP))

	baz := pkg.File("baz")
	bazWP := &writerPrepperNamer{
		str:  `var test = "testing update namer 2"`,
		name: "test",
	}
	assert.NoError(t, baz.AddWriterTo(bazWP))
}

func TestNamerRenameCausesCollision(t *testing.T) {
	ctx := NewMemoryContext()
	pkg := ctx.MustPackage("foo")

	bar := pkg.File("bar")
	barWP := &writerPrepperNamer{
		str:  `var ok = "testing name collision 1"`,
		name: "ok",
	}
	assert.NoError(t, bar.AddWriterTo(barWP))

	baz := pkg.File("baz")
	bazWP := &writerPrepperNamer{
		str:  `var test = "testing name collision 2"`,
		name: "test",
	}
	assert.NoError(t, baz.AddWriterTo(bazWP))

	barWP.name = "test"
	barWP.str = `var test = "testing name collision 1.1"`
	err := bar.UpdateNamer(barWP)
	assert.Equal(t, "UpdateNamer in file bar: Name 'test' already exists in package 'foo'", err.Error())
}
