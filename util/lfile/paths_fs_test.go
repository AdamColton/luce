package lfile_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/adamcolton/luce/util/lfile"
	"github.com/stretchr/testify/assert"
)

func TestPathsFS(t *testing.T) {
	files := fstest.MapFS{
		"a.txt":     {Data: []byte("aaa")},
		"dir/b.txt": {Data: []byte("bb")},
	}
	src := lfile.Paths{"a.txt", "dir/b.txt"}.FS(files)
	var _ lfile.IteratorSource = src

	i, done := src.Iterator()
	assert.False(t, done)
	path, done := i.Cur()
	assert.Equal(t, "a.txt", path)
	assert.False(t, done)
	assert.Equal(t, 0, i.Idx())
	assert.Equal(t, []byte("aaa"), i.Data())
	assert.Equal(t, int64(3), i.Stat().Size())
	assert.NoError(t, i.Err())

	// Next gives the path of the file it moved to
	path, done = i.Next()
	assert.Equal(t, "dir/b.txt", path)
	assert.False(t, done)
	assert.Equal(t, 1, i.Idx())
	assert.Equal(t, "dir/b.txt", i.Path())
	assert.Equal(t, []byte("bb"), i.Data())
	assert.Equal(t, int64(2), i.Stat().Size())

	path, done = i.Next()
	assert.Equal(t, "", path)
	assert.True(t, done)
	assert.True(t, i.Done())
	assert.Equal(t, "", i.Path())

	// Reset starts again
	assert.False(t, i.Reset())
	assert.Equal(t, "a.txt", i.Path())
	assert.Equal(t, []byte("aaa"), i.Data())
}

func TestPathsFSErrors(t *testing.T) {
	files := fstest.MapFS{"a.txt": {Data: []byte("aaa")}}

	// a file that can't be read stops the iteration
	i, done := lfile.Paths{"missing", "a.txt"}.FS(files).Iterator()
	assert.False(t, done)
	assert.Nil(t, i.Data())
	assert.ErrorIs(t, i.Err(), fs.ErrNotExist)
	assert.True(t, i.Done())
	assert.Nil(t, i.Data(), "the error is kept")

	// so does one that can't be described, at the next file
	i, done = lfile.Paths{"missing", "a.txt"}.FS(files).Iterator()
	assert.False(t, done)
	assert.Nil(t, i.Stat())
	assert.ErrorIs(t, i.Err(), fs.ErrNotExist)
	assert.False(t, i.Done())
	_, done = i.Next()
	assert.True(t, done)

	// no paths
	i, done = lfile.Paths(nil).FS(files).Iterator()
	assert.True(t, done)
	assert.True(t, i.Done())
	assert.Equal(t, "", i.Path())
}

func TestPathsOS(t *testing.T) {
	dir := t.TempDir()
	name := filepath.ToSlash(filepath.Join(dir, "a.txt"))
	assert.NoError(t, os.WriteFile(name, []byte("hello"), 0o644))

	// Paths and a PathsFS with no CoreFS read from the OS
	for _, src := range []lfile.IteratorSource{
		lfile.Paths{name},
		lfile.PathsFS{Paths: lfile.Paths{name}},
	} {
		i, done := src.Iterator()
		assert.False(t, done)
		assert.Equal(t, []byte("hello"), i.Data())
		assert.Equal(t, int64(5), i.Stat().Size())
		assert.NoError(t, i.Err())
	}
}
