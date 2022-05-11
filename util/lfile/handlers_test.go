package lfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiHandler(t *testing.T) {
	restore := setupForTestMultiHandler()
	defer restore()

	fs := Paths{"foo/", "foo.txt", "bar.txt", "bar/"}
	c := GetContentsHandler{}
	bt := &GetByTypeHandler{}
	err := RunHandlerSource(fs, MultiHandler{c, bt})
	assert.NoError(t, err)

	expected := GetContentsHandler{
		"foo.txt": []byte("foo.txt"),
		"bar.txt": []byte("bar.txt"),
	}
	assert.Equal(t, expected, c)
	assert.Equal(t, []string{"foo.txt", "bar.txt"}, bt.Files)
	assert.Equal(t, []string{"foo/", "bar/"}, bt.Dirs)
}

func setupForTestMultiHandler() func() {
	restoreStat := Stat
	restoreReadFile := ReadFile
	Stat = mockStat
	ReadFile = mockReadFileAsName
	return func() {
		ReadFile = restoreReadFile
		Stat = restoreStat
	}
}
