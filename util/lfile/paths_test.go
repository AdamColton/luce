package lfile

import (
	"io/fs"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func mockReadFileAsName(filename string) ([]byte, error) {
	return []byte(filename), nil
}

type mockFileInfo struct {
	isDir bool
}

func (mfi mockFileInfo) Name() string {
	return ""
}
func (mfi mockFileInfo) Size() int64 {
	return 0
}
func (mfi mockFileInfo) Mode() fs.FileMode {
	return 0
}
func (mfi mockFileInfo) ModTime() time.Time {
	return time.Now()
}
func (mfi mockFileInfo) IsDir() bool {
	return mfi.isDir
}
func (mfi mockFileInfo) Sys() interface{} {
	return nil
}

func mockStat(name string) (os.FileInfo, error) {
	return mockFileInfo{
		isDir: strings.HasSuffix(name, "/"),
	}, nil
}

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
	restore := setupForTestIter()
	defer restore()

	fs := Paths{"foo.txt", "bar.txt"}
	c := 0
	i, done := fs.Iterator()
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

	mh := &mockHandler{
		autoload: false,
	}
	err := RunHandlerSource(fs, mh)
	assert.NoError(t, err)
	assert.Equal(t, []string(fs), mh.got)

	mh.got = nil
	err = RunHandler(i, mh)
	assert.NoError(t, err)
	assert.Equal(t, []string(fs), mh.got)

}

func setupForTestIter() func() {
	restoreReadFile := ReadFile
	restoreStat := Stat
	ReadFile = mockReadFileAsName
	Stat = mockStat
	return func() {
		ReadFile = restoreReadFile
		Stat = restoreStat
	}
}
