package lfile_test

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/adamcolton/luce/util/lfile"
	"github.com/stretchr/testify/assert"
)

func TestOSRepository(t *testing.T) {
	// the paths are slash separated, whatever the operating system
	dir := filepath.ToSlash(t.TempDir())
	var repo lfile.Repository = lfile.OSRepository{}

	f, err := repo.Create(dir + "/a.txt")
	assert.NoError(t, err)
	_, err = f.(io.Writer).Write([]byte("hello"))
	assert.NoError(t, err)
	assert.NoError(t, f.Close())

	f, err = repo.Open(dir + "/a.txt")
	assert.NoError(t, err)
	data, err := io.ReadAll(f)
	assert.NoError(t, err)
	assert.Equal(t, "hello", string(data))
	assert.NoError(t, f.Close())

	data, err = repo.ReadFile(dir + "/a.txt")
	assert.NoError(t, err)
	assert.Equal(t, "hello", string(data))

	assert.NoError(t, os.Symlink(dir+"/a.txt", dir+"/link"))
	info, err := repo.Stat(dir + "/link")
	assert.NoError(t, err)
	assert.Equal(t, int64(5), info.Size(), "Stat follows the link")
	assert.True(t, info.Mode().IsRegular())
	info, err = repo.Lstat(dir + "/link")
	assert.NoError(t, err)
	assert.Equal(t, fs.ModeSymlink, info.Mode()&fs.ModeSymlink, "Lstat does not")

	entries, err := repo.ReadDir(dir)
	assert.NoError(t, err)
	if assert.Len(t, entries, 2) {
		assert.Equal(t, "a.txt", entries[0].Name())
		assert.Equal(t, "link", entries[1].Name())
	}

	assert.NoError(t, repo.Remove(dir+"/link"))
	_, err = repo.Lstat(dir + "/link")
	assert.ErrorIs(t, err, fs.ErrNotExist)
}

func TestOSRepositoryErrors(t *testing.T) {
	dir := filepath.ToSlash(t.TempDir())
	var repo lfile.Repository = lfile.OSRepository{}
	missing := dir + "/missing"

	f, err := repo.Open(missing)
	assert.ErrorIs(t, err, fs.ErrNotExist)
	assert.True(t, f == nil, "a failed Open returns a nil fs.File, not a nil *os.File in one")

	f, err = repo.Create(missing + "/a.txt")
	assert.ErrorIs(t, err, fs.ErrNotExist)
	assert.True(t, f == nil)

	assert.ErrorIs(t, repo.Remove(missing), fs.ErrNotExist)
	_, err = repo.Stat(missing)
	assert.ErrorIs(t, err, fs.ErrNotExist)
	_, err = repo.Lstat(missing)
	assert.ErrorIs(t, err, fs.ErrNotExist)
	_, err = repo.ReadFile(missing)
	assert.ErrorIs(t, err, fs.ErrNotExist)
	_, err = repo.ReadDir(missing)
	assert.ErrorIs(t, err, fs.ErrNotExist)
}
