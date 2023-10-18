package lfile_test

import (
	"io/fs"
	"os"
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/lfile"
	"github.com/adamcolton/luce/util/lfile/lfilemock"
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
	// mock := mockDir{
	// 	name: "foo/bar/baz/",
	// 	entries: []fs.DirEntry{
	// 		mockDirEntry{
	// 			name:  "dir2",
	// 			isDir: true,
	// 		},
	// 		mockDirEntry{
	// 			name:  "x.txt",
	// 			isDir: false,
	// 		},
	// 		mockDirEntry{
	// 			name:  "dir1",
	// 			isDir: true,
	// 		},
	// 		mockDirEntry{
	// 			name:  "b.txt",
	// 			isDir: false,
	// 		},
	// 	},
	// }

	mock := lfilemock.Parse(map[string]any{
		"foo": map[string]any{
			"bar": map[string]any{
				"baz": map[string]any{
					"dir2":  map[string]any{},
					"x.txt": "",
					"dir1":  map[string]any{},
					"b.txt": "",
				},
			},
		},
	})

	// lfile.GetDirContents is looking at repo.Name()
	// this is set in node.File()
	node, found := mock.Get("foo/bar/baz")
	assert.True(t, found)
	repo := node.File()

	got, err := lfile.GetDirContents(repo)
	assert.NoError(t, err)

	expected := &lfile.DirContents{
		Name:    "baz",
		Path:    "foo/bar/",
		SubDirs: []string{"dir1", "dir2"},
		Files:   []string{"b.txt", "x.txt"},
	}
	assert.Equal(t, expected, got)

	repo.Err = lerr.Str("test error")
	got, err = lfile.GetDirContents(repo)
	assert.Equal(t, repo.Err, err)
	assert.Nil(t, got)
}

func TestFoo(t *testing.T) {
	var f fs.File
	var err error
	f, err = os.Open("lfilemock")
	assert.NoError(t, err)

	s, err := f.Stat()
	assert.NoError(t, err)

	assert.True(t, s.IsDir())
}
