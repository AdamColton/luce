package lfile_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
	"time"

	"github.com/adamcolton/luce/util/lfile"
	"github.com/adamcolton/luce/util/lfile/lfilemock"
	"github.com/stretchr/testify/assert"
)

func TestRepositorySize(t *testing.T) {
	r := lfilemock.Parse(map[string]any{
		"file1.txt": "this is test file 1",
		"dir1": map[string]any{
			"file2.bin": []byte{3, 1, 4, 1, 5, 9, 2, 6, 5},
		},
		"dir2": map[string]any{
			"dir3":      map[string]any{},
			"file4.txt": "this is file 4",
		},
	}).Repository()

	size, err := lfile.Size(r, "/")
	assert.NoError(t, err)
	assert.Equal(t, int64(42), size)
}

func TestSizeMapFS(t *testing.T) {
	m := fstest.MapFS{
		"a.txt":         {Data: []byte("12345")},
		"docs/b.txt":    {Data: []byte("123")},
		"docs/sub/c.md": {Data: []byte("1234567")},
		"empty/x":       {Data: nil},
	}

	size, err := lfile.Size(m, ".")
	assert.NoError(t, err)
	assert.Equal(t, int64(15), size)

	size, err = lfile.Size(m, "docs")
	assert.NoError(t, err)
	assert.Equal(t, int64(10), size)

	size, err = lfile.Size(m, "a.txt")
	assert.NoError(t, err)
	assert.Equal(t, int64(5), size)

	_, err = lfile.Size(m, "missing")
	assert.ErrorIs(t, err, fs.ErrNotExist)
}

// listedFS lists a file in its directory that can't be opened, and fails to
// list the directory "bad".
type listedFS struct {
	fstest.MapFS
}

func (l listedFS) ReadDir(name string) ([]fs.DirEntry, error) {
	if name == "bad" {
		return nil, fs.ErrPermission
	}
	des, err := l.MapFS.ReadDir(name)
	if name == "ghosts" {
		// first, so the entries after it are skipped once it has failed
		ghost := fs.FileInfoToDirEntry(fakeInfo{name: "a-ghost"})
		des = append([]fs.DirEntry{ghost}, des...)
	}
	return des, err
}

type fakeInfo struct {
	name string
	mode fs.FileMode
}

func (f fakeInfo) Name() string       { return f.name }
func (f fakeInfo) Size() int64        { return 1000 }
func (f fakeInfo) Mode() fs.FileMode  { return f.mode }
func (f fakeInfo) ModTime() time.Time { return time.Time{} }
func (f fakeInfo) IsDir() bool        { return f.mode.IsDir() }
func (f fakeInfo) Sys() any           { return nil }

func TestSizeErrors(t *testing.T) {
	l := listedFS{fstest.MapFS{
		"ok/a.txt":     {Data: []byte("a")},
		"bad/x":        {Data: []byte("x")},
		"ghosts/a.txt": {Data: []byte("a")},
	}}

	// a directory that can't be listed, at any depth
	_, err := lfile.Size(l, "bad")
	assert.ErrorIs(t, err, fs.ErrPermission)
	_, err = lfile.Size(l, ".")
	assert.Error(t, err)

	// listed but gone by the time it is described
	_, err = lfile.Size(l, "ghosts")
	assert.ErrorIs(t, err, fs.ErrNotExist)

	// the good part of the tree is fine
	size, err := lfile.Size(l, "ok")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), size)
}

// linkFS reports "link" as a symbolic link to a 1000 byte file.
type linkFS struct {
	fstest.MapFS
}

func (l linkFS) Lstat(name string) (fs.FileInfo, error) {
	if name == "link" || name == "dir/link" {
		return fakeInfo{name: name, mode: fs.ModeSymlink}, nil
	}
	return l.MapFS.Stat(name)
}

func TestSizeSkipsLinks(t *testing.T) {
	l := linkFS{fstest.MapFS{
		"link":     {Data: []byte("1234567890")},
		"dir/link": {Data: []byte("1234567890")},
		"dir/a":    {Data: []byte("123")},
	}}
	size, err := lfile.Size(l, ".")
	assert.NoError(t, err)
	assert.Equal(t, int64(3), size)
}

func TestSizeOS(t *testing.T) {
	dir := t.TempDir()
	assert.NoError(t, os.Mkdir(filepath.Join(dir, "sub"), 0o755))
	assert.NoError(t, os.WriteFile(filepath.Join(dir, "a.txt"), []byte("12345"), 0o644))
	assert.NoError(t, os.WriteFile(filepath.Join(dir, "sub", "b.txt"), []byte("123"), 0o644))
	// a link to a directory is not followed
	assert.NoError(t, os.Symlink(dir, filepath.Join(dir, "sub", "loop")))

	size, err := lfile.Size(lfile.OSRepository{}, dir)
	assert.NoError(t, err)
	assert.Equal(t, int64(8), size)
}

// lstatFile has its own Lstat.
type lstatFile struct {
	fs.File
	info fs.FileInfo
}

func (l lstatFile) Lstat() (fs.FileInfo, error) { return l.info, nil }

type lstatFileFS struct {
	coreOnly
}

func (l lstatFileFS) Open(name string) (fs.File, error) {
	f, err := l.coreOnly.Open(name)
	if err != nil {
		return nil, err
	}
	return lstatFile{File: f, info: fakeInfo{name: "link", mode: fs.ModeSymlink}}, nil
}

func TestWrapLstat(t *testing.T) {
	m := fstest.MapFS{"a.txt": {Data: []byte("12345")}}

	t.Run("FSLstater", func(t *testing.T) {
		info, err := lfile.WrapLstat(linkFS{m})("link")
		assert.NoError(t, err)
		assert.Equal(t, fs.ModeSymlink, info.Mode())
	})

	t.Run("FileLstater", func(t *testing.T) {
		core := lstatFileFS{newCoreOnly(m)}
		info, err := lfile.WrapLstat(core)("a.txt")
		assert.NoError(t, err)
		assert.Equal(t, fs.ModeSymlink, info.Mode())
		assert.Equal(t, 1, *core.closed)
	})

	t.Run("Stat", func(t *testing.T) {
		core := newCoreOnly(m)
		lstat := lfile.WrapLstat(core)
		info, err := lstat("a.txt")
		assert.NoError(t, err)
		assert.Equal(t, int64(5), info.Size())
		assert.Equal(t, 1, *core.closed, "the file is closed")

		_, err = lstat("missing")
		assert.ErrorIs(t, err, fs.ErrNotExist)
	})
}
