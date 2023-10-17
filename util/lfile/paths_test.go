package lfile

import (
	"io/fs"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// nameFS is a CoreFS in which every file holds its own name, and a name that
// ends in a slash is a directory.
type nameFS struct{}

func (nameFS) Open(name string) (fs.File, error)          { return nil, fs.ErrNotExist }
func (nameFS) ReadDir(name string) ([]fs.DirEntry, error) { return nil, fs.ErrNotExist }
func (nameFS) ReadFile(name string) ([]byte, error)       { return []byte(name), nil }
func (nameFS) Stat(name string) (fs.FileInfo, error) {
	return nameInfo{name: name, isDir: name != "" && name[len(name)-1] == '/'}, nil
}

type nameInfo struct {
	name  string
	isDir bool
}

func (ni nameInfo) Name() string       { return ni.name }
func (ni nameInfo) Size() int64        { return int64(len(ni.name)) }
func (ni nameInfo) Mode() fs.FileMode  { return 0 }
func (ni nameInfo) ModTime() time.Time { return time.Time{} }
func (ni nameInfo) IsDir() bool        { return ni.isDir }
func (ni nameInfo) Sys() any           { return nil }

type mockHandler struct {
	got      []string
	autoload bool
}

func (mh *mockHandler) HandleIter(i Iterator) {
	mh.got = append(mh.got, i.Path())
}

func (mh *mockHandler) Autoload() bool {
	return mh.autoload
}

func TestIter(t *testing.T) {
	fs := Paths{"foo.txt", "bar.txt"}

	i, done := fs.FS(nameFS{}).Iterator()
	c := 0
	for ; !done; _, done = i.Next() {
		c++
		assert.Equal(t, fs[i.(*pathsIterator).Index], string(i.Data()))
		assert.False(t, i.Done())
		assert.False(t, i.Stat().IsDir())
	}
	assert.True(t, i.Done())
	assert.Equal(t, len(fs), c)

	c = 0
	for done = i.Reset(); !done; _, done = i.Next() {
		c++
		assert.Equal(t, fs[i.(*pathsIterator).Index], string(i.Data()))
		assert.False(t, i.Done())
		assert.False(t, i.Stat().IsDir())
	}
	assert.True(t, i.Done())
	assert.Equal(t, len(fs), c)

	// mh := &mockHandler{
	// 	autoload: false,
	// }
	// err := RunHandlerSource(fs, mh)
	// assert.NoError(t, err)
	// assert.Equal(t, []string(fs), mh.got)

	// mh.got = nil
	// err = RunHandler(i, mh)
	// assert.NoError(t, err)
	// assert.Equal(t, []string(fs), mh.got)

}
