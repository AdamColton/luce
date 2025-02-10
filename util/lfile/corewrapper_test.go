package lfile_test

import (
	"io"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/adamcolton/luce/util/lfile"
	"github.com/adamcolton/luce/util/upgrade"
	"github.com/stretchr/testify/assert"
)

// coreOnly is a CoreFS and nothing more, like embed.FS. It counts the files it
// has opened and closed.
type coreOnly struct {
	m      fstest.MapFS
	opened *int
	closed *int
}

func newCoreOnly(m fstest.MapFS) coreOnly {
	return coreOnly{m: m, opened: new(int), closed: new(int)}
}

func (c coreOnly) Open(name string) (fs.File, error) {
	f, err := c.m.Open(name)
	if err != nil {
		return nil, err
	}
	*c.opened++
	return countClose{File: f, closed: c.closed}, nil
}

func (c coreOnly) ReadFile(name string) ([]byte, error)       { return c.m.ReadFile(name) }
func (c coreOnly) ReadDir(name string) ([]fs.DirEntry, error) { return c.m.ReadDir(name) }

type countClose struct {
	fs.File
	closed *int
}

func (c countClose) Close() error {
	*c.closed++
	return c.File.Close()
}

func TestWrapCoreFSStat(t *testing.T) {
	core := newCoreOnly(fstest.MapFS{
		"docs/a.txt": {Data: []byte("hello")},
	})
	var _ lfile.FSStater = core.m // MapFS has Stat, coreOnly does not
	_, hasStat := any(core).(lfile.FSStater)
	assert.False(t, hasStat)

	fsr := lfile.WrapCoreFS(core)

	info, err := fsr.Stat("docs/a.txt")
	assert.NoError(t, err)
	assert.Equal(t, "a.txt", info.Name())
	assert.Equal(t, int64(5), info.Size())
	assert.False(t, info.IsDir())
	assert.Equal(t, 1, *core.opened)
	assert.Equal(t, 1, *core.closed, "Stat closes the file it opens")

	info, err = fsr.Stat("docs")
	assert.NoError(t, err)
	assert.True(t, info.IsDir())

	_, err = fsr.Stat("missing")
	assert.ErrorIs(t, err, fs.ErrNotExist)
	assert.Equal(t, *core.opened, *core.closed)
}

func TestWrapCoreFSUsesStatOfWrapped(t *testing.T) {
	m := fstest.MapFS{"a.txt": {Data: []byte("hi")}}
	info, err := lfile.WrapCoreFS(m).Stat("a.txt")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), info.Size())
}

func TestWrapCoreFSWrapped(t *testing.T) {
	core := newCoreOnly(fstest.MapFS{})
	fsr := lfile.WrapCoreFS(core)

	w, ok := fsr.(interface{ Wrapped() any })
	assert.True(t, ok)
	assert.Equal(t, core.opened, w.Wrapped().(coreOnly).opened)

	// upgrade.To looks through the wrapper
	m := fstest.MapFS{}
	_, ok = upgrade.To[fstest.MapFS](lfile.WrapCoreFS(m))
	assert.True(t, ok)
}

func TestWrapCoreFSOpen(t *testing.T) {
	core := newCoreOnly(fstest.MapFS{
		"dir/a.txt": {Data: []byte("a")},
		"dir/b.txt": {Data: []byte("b")},
		"dir/c.txt": {Data: []byte("c")},
	})
	fsr := lfile.WrapCoreFS(core)

	_, err := fsr.Open("nope")
	assert.ErrorIs(t, err, fs.ErrNotExist)

	f, err := fsr.Open("dir")
	assert.NoError(t, err)
	dir, ok := f.(lfile.Dir)
	assert.True(t, ok, "files from a wrapped CoreFS are Dirs")
	assert.Equal(t, "dir", dir.Name())

	// read in pages, as os.File does
	des, err := dir.ReadDir(2)
	assert.NoError(t, err)
	if assert.Len(t, des, 2) {
		assert.Equal(t, "a.txt", des[0].Name())
		assert.Equal(t, "b.txt", des[1].Name())
	}
	names, err := dir.Readdirnames(2)
	assert.NoError(t, err)
	assert.Equal(t, []string{"c.txt"}, names)
	des, err = dir.ReadDir(2)
	assert.Nil(t, des)
	assert.Equal(t, io.EOF, err)
	names, err = dir.Readdirnames(2)
	assert.Nil(t, names)
	assert.Equal(t, io.EOF, err)

	// n <= 0 reads whatever is left, without an error
	des, err = dir.ReadDir(-1)
	assert.NoError(t, err)
	assert.Empty(t, des)
	assert.NoError(t, f.Close())

	f, err = fsr.Open("dir")
	assert.NoError(t, err)
	names, err = f.(lfile.Dir).Readdirnames(-1)
	assert.NoError(t, err)
	assert.Equal(t, []string{"a.txt", "b.txt", "c.txt"}, names)
	f.Close()
	assert.Equal(t, *core.opened, *core.closed)
}

func TestWrapCoreFSReadDirErrors(t *testing.T) {
	fsr := lfile.WrapCoreFS(newCoreOnly(fstest.MapFS{
		"a.txt": {Data: []byte("a")},
	}))

	f, err := fsr.Open("a.txt")
	assert.NoError(t, err)
	_, err = f.(lfile.Dir).ReadDir(-1)
	assert.Error(t, err, "a file can't be listed")
	_, err = f.(lfile.Dir).Readdirnames(-1)
	assert.Error(t, err)
}
