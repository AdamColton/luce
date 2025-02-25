package lfile_test

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
	"testing/fstest"

	"github.com/adamcolton/luce/util/lfile"
	"github.com/stretchr/testify/assert"
)

var filesTree = fstest.MapFS{
	"books/a.mp3":        {},
	"books/sub/b.mp3":    {},
	"books/sub/.x/c.mp3": {},
	"books/notes.txt":    {},
	"other.mp3":          {},
}

func sorted(files []string) []string {
	sort.Strings(files)
	return files
}

func TestRootGetFiles(t *testing.T) {
	recursive := lfile.MatchExt(true, "mp3")

	// the paths do not start with the root, however it is written
	for _, root := range []string{"books", "books/", "./books", "books/."} {
		files, err := lfile.RootGetFiles(root, filesTree, recursive)
		assert.NoError(t, err)
		assert.Equal(t, []string{"a.mp3", "sub/b.mp3"}, sorted(files), root)
	}

	// the current directory has no prefix to remove
	for _, root := range []string{"", ".", "./"} {
		files, err := lfile.RootGetFiles(root, filesTree, recursive)
		assert.NoError(t, err)
		assert.Equal(t, []string{"books/a.mp3", "books/sub/b.mp3", "other.mp3"}, sorted(files), root)
	}

	// not recursive is only the directory itself
	files, err := lfile.RootGetFiles("books", filesTree, lfile.MatchExt(false, "mp3"))
	assert.NoError(t, err)
	assert.Equal(t, []string{"a.mp3"}, []string(files))
	files, err = lfile.RootGetFiles("books", filesTree, lfile.MatchExt(true, "mp3", "txt"))
	assert.NoError(t, err)
	assert.Equal(t, []string{"a.mp3", "notes.txt", "sub/b.mp3"}, sorted(files))

	// an error ends the walk
	_, err = lfile.RootGetFiles("nothing", filesTree, recursive)
	assert.Error(t, err)
}

func TestRootGetFilesOS(t *testing.T) {
	dir := t.TempDir()
	assert.NoError(t, os.MkdirAll(filepath.Join(dir, "sub"), 0o755))
	assert.NoError(t, os.WriteFile(filepath.Join(dir, "a.mp3"), nil, 0o644))
	assert.NoError(t, os.WriteFile(filepath.Join(dir, "sub", "b.mp3"), nil, 0o644))

	// a nil file system is the operating system
	files, err := lfile.RootGetFiles(filepath.ToSlash(dir), nil, lfile.MatchExt(true, "mp3"))
	assert.NoError(t, err)
	assert.Equal(t, []string{"a.mp3", "sub/b.mp3"}, sorted(files))
}

func TestFilesByExtAndRegex(t *testing.T) {
	files, err := lfile.FilesByExt("books", filesTree, "mp3", "txt")
	assert.NoError(t, err)
	assert.Equal(t, []string{"a.mp3", "notes.txt", "sub/b.mp3"}, sorted(files))

	files, err = lfile.FilesByRegex("books", filesTree, `notes\.txt$`)
	assert.NoError(t, err)
	assert.Equal(t, []string{"notes.txt"}, []string(files))

	_, err = lfile.FilesByRegex("books", filesTree, "(")
	assert.Error(t, err)
}
