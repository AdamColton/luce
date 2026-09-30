package lfile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/lfile"
	"github.com/adamcolton/luce/util/lfile/lfilemock"
	"github.com/stretchr/testify/assert"
)

func TestTryRemoveOS(t *testing.T) {
	dir := filepath.ToSlash(t.TempDir())
	var r lfile.FSRemover = lfile.OSRepository{}

	assert.NoError(t, os.WriteFile(dir+"/a.txt", nil, 0o644))
	assert.NoError(t, lfile.TryRemove(r, dir+"/a.txt"))
	_, err := os.Stat(dir + "/a.txt")
	assert.ErrorIs(t, err, os.ErrNotExist)

	// what is not there, or is below what is not there, is not an error
	assert.NoError(t, lfile.TryRemove(r, dir+"/a.txt"))
	assert.NoError(t, lfile.TryRemove(r, dir+"/missing/b.txt"))

	// but a directory that has something in it is
	assert.NoError(t, os.Mkdir(dir+"/full", 0o755))
	assert.NoError(t, os.WriteFile(dir+"/full/x", nil, 0o644))
	assert.Error(t, lfile.TryRemove(r, dir+"/full"))
	assert.NoError(t, os.Remove(dir+"/full/x"))
	assert.NoError(t, lfile.TryRemove(r, dir+"/full"))

	// and so is a path through a file
	assert.NoError(t, os.WriteFile(dir+"/file", nil, 0o644))
	assert.Error(t, lfile.TryRemove(r, dir+"/file/x"))
}

func TestTryRemoveMock(t *testing.T) {
	root := lfilemock.Parse(map[string]any{
		"a.txt": "a",
		"full":  map[string]any{"x": "x"},
	})
	repo := root.Repository()
	r := repo.(lfile.FSRemover)

	assert.NoError(t, lfile.TryRemove(r, "a.txt"))
	_, found := root.Get("a.txt")
	assert.False(t, found)
	assert.NoError(t, lfile.TryRemove(r, "a.txt"))
	assert.NoError(t, lfile.TryRemove(r, "missing/b.txt"))
	assert.Error(t, lfile.TryRemove(r, "full"))

	// any other error is returned
	root.Err = lerr.Str("boom")
	assert.Equal(t, lerr.Str("boom"), lfile.TryRemove(r, "a.txt"))
}
