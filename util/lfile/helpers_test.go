package lfile_test

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/adamcolton/luce/util/lfile"
	"github.com/stretchr/testify/assert"
)

func TestCoreFsStat(t *testing.T) {
	m := fstest.MapFS{"docs/a.txt": {Data: []byte("hello")}}

	// a file system with Stat is asked directly
	info, err := lfile.CoreFsStat(m, "docs/a.txt")
	assert.NoError(t, err)
	assert.Equal(t, int64(5), info.Size())

	// a file system without one is opened and closed
	core := newCoreOnly(m)
	info, err = lfile.CoreFsStat(core, "docs/a.txt")
	assert.NoError(t, err)
	assert.Equal(t, "a.txt", info.Name())
	assert.Equal(t, 1, *core.opened)
	assert.Equal(t, 1, *core.closed)

	_, err = lfile.CoreFsStat(core, "missing")
	assert.ErrorIs(t, err, fs.ErrNotExist)
}

// namesFS opens files that list their own names.
type namesFS struct {
	coreOnly
	names []string
	err   error
}

func (n namesFS) Open(name string) (fs.File, error) {
	f, err := n.coreOnly.Open(name)
	if err != nil {
		return nil, err
	}
	return namesFile{File: f, names: n.names, err: n.err}, nil
}

type namesFile struct {
	fs.File
	names []string
	err   error
}

func (n namesFile) Readdirnames(int) ([]string, error) { return n.names, n.err }

func TestReadDirNames(t *testing.T) {
	m := fstest.MapFS{
		"dir/b.txt": {Data: []byte("b")},
		"dir/a.txt": {Data: []byte("a")},
		"dir/sub/c": {Data: []byte("c")},
	}
	want := []string{"a.txt", "b.txt", "sub"}

	t.Run("Fallback", func(t *testing.T) {
		// the files of a MapFS can't list names, so the directory is read
		names, err := lfile.ReadDirNames(m, "dir")
		assert.NoError(t, err)
		assert.Equal(t, want, names)

		core := newCoreOnly(m)
		names, err = lfile.ReadDirNames(core, "dir")
		assert.NoError(t, err)
		assert.Equal(t, want, names)
		assert.Equal(t, *core.opened, *core.closed, "no file is left open")

		_, err = lfile.ReadDirNames(m, "dir/a.txt")
		assert.Error(t, err, "not a directory")
	})

	t.Run("WrappedCoreFS", func(t *testing.T) {
		core := newCoreOnly(m)
		names, err := lfile.ReadDirNames(lfile.WrapCoreFS(core), "dir")
		assert.NoError(t, err)
		assert.Equal(t, want, names)
		assert.Equal(t, *core.opened, *core.closed, "no file is left open")
	})

	t.Run("OS", func(t *testing.T) {
		dir := t.TempDir()
		assert.NoError(t, os.WriteFile(filepath.Join(dir, "b"), nil, 0o644))
		assert.NoError(t, os.WriteFile(filepath.Join(dir, "a"), nil, 0o644))

		names, err := lfile.ReadDirNames(lfile.OSRepository{}, dir)
		assert.NoError(t, err)
		assert.Equal(t, []string{"a", "b"}, names)

		_, err = lfile.ReadDirNames(lfile.OSRepository{}, filepath.Join(dir, "a"))
		assert.Error(t, err)
	})

	t.Run("Errors", func(t *testing.T) {
		_, err := lfile.ReadDirNames(m, "missing")
		assert.ErrorIs(t, err, fs.ErrNotExist)

		boom := fs.ErrPermission
		core := namesFS{coreOnly: newCoreOnly(m), names: []string{"z", "y"}, err: boom}
		names, err := lfile.ReadDirNames(core, "dir")
		assert.Nil(t, names)
		assert.ErrorIs(t, err, boom)
		assert.Equal(t, 1, *core.closed)

		// io.EOF is not a failure, and the names are sorted
		core.err = io.EOF
		names, err = lfile.ReadDirNames(core, "dir")
		assert.NoError(t, err)
		assert.Equal(t, []string{"y", "z"}, names)
	})
}
