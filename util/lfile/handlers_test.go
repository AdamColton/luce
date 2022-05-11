package lfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNamesByType(t *testing.T) {
	restore := setupForTestGetNamesByType()
	defer restore()

	fs := Paths{"foo/", "foo.txt", "bar.txt", "bar/"}
	bt := &GetByTypeHandler{}
	err := RunHandlerSource(fs, bt)
	assert.NoError(t, err)
	assert.Equal(t, []string{"foo.txt", "bar.txt"}, bt.Files)
	assert.Equal(t, []string{"foo/", "bar/"}, bt.Dirs)
}

func setupForTestGetNamesByType() func() {
	restoreStat := Stat
	Stat = mockStat
	return func() {
		Stat = restoreStat
	}
}

func TestGetContents(t *testing.T) {
	restore := setupForTestGetContents()
	defer restore()

	fs := Paths{"foo/", "foo.txt", "bar.txt", "bar/"}
	c := GetContentsHandler{}
	err := RunHandlerSource(fs, c)
	assert.NoError(t, err)

	expected := GetContentsHandler{
		"foo.txt": []byte("foo.txt"),
		"bar.txt": []byte("bar.txt"),
	}
	assert.Equal(t, expected, c)
}

func setupForTestGetContents() func() {
	restoreStat := Stat
	restoreReadFile := ReadFile
	Stat = mockStat
	ReadFile = mockReadFileAsName
	return func() {
		ReadFile = restoreReadFile
		Stat = restoreStat
	}
}
