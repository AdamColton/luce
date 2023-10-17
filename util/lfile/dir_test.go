package lfile

import (
	"io/fs"
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/stretchr/testify/assert"
)

type mockDir struct {
	name    string
	entries []fs.DirEntry
	err     error
}

func (md mockDir) ReadDir(n int) ([]fs.DirEntry, error) {
	if md.err != nil {
		return nil, md.err
	}
	ln := len(md.entries)
	if n < 0 {
		n = ln
	} else if n > ln {
		n = ln
	}
	return md.entries[:ln], nil
}

func (md mockDir) Name() string {
	return md.name
}

type mockDirEntry struct {
	name  string
	isDir bool
}

func (md mockDirEntry) Name() string {
	return md.name
}
func (md mockDirEntry) IsDir() bool {
	return md.isDir
}
func (md mockDirEntry) Type() fs.FileMode {
	return 0
}
func (md mockDirEntry) Info() (fs.FileInfo, error) {
	return nil, nil
}

func TestGetDirContents(t *testing.T) {
	mock := mockDir{
		name: "foo/bar/baz/",
		entries: []fs.DirEntry{
			mockDirEntry{
				name:  "dir2",
				isDir: true,
			},
			mockDirEntry{
				name:  "x.txt",
				isDir: false,
			},
			mockDirEntry{
				name:  "dir1",
				isDir: true,
			},
			mockDirEntry{
				name:  "b.txt",
				isDir: false,
			},
		},
	}

	got, err := GetDirContents(mock)
	assert.NoError(t, err)

	expected := &DirContents{
		Name:    "baz",
		Path:    "foo/bar/",
		SubDirs: []string{"dir1", "dir2"},
		Files:   []string{"b.txt", "x.txt"},
	}
	assert.Equal(t, expected, got)

	mock.err = lerr.Str("test error")
	got, err = GetDirContents(mock)
	assert.Equal(t, mock.err, err)
	assert.Nil(t, got)
}
