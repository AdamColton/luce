package lfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetByTypeHandler(t *testing.T) {
	restore := setupForHandlersTest()
	defer restore()

	fs := Paths{"foo/", "foo.txt", "bar.txt", "bar/"}
	bt := &GetByTypeHandler{}
	err := RunHandlerSource(fs, bt)
	assert.NoError(t, err)

	expected := &GetByTypeHandler{
		Files: []string{"foo.txt", "bar.txt"},
		Dirs:  []string{"foo/", "bar/"},
	}
	assert.Equal(t, expected, bt)
}

func TestGetContentsHandler(t *testing.T) {
	restore := setupForHandlersTest()
	defer restore()

	fs := Paths{"foo/", "foo.txt", "bar.txt", "bar/"}
	gt := make(GetContentsHandler)
	err := RunHandlerSource(fs, gt)
	assert.NoError(t, err)

	expected := GetContentsHandler{
		"foo.txt": []byte("foo.txt"),
		"bar.txt": []byte("bar.txt"),
	}
	assert.Equal(t, expected, gt)
}

func TestMultiHandler(t *testing.T) {
	restore := setupForHandlersTest()
	defer restore()

	fs := Paths{"foo/", "foo.txt", "bar.txt", "bar/"}
	c := GetContentsHandler{}
	bt := &GetByTypeHandler{}
	err := RunHandlerSource(fs, MultiHandler{c, bt})
	assert.NoError(t, err)

	expectedGC := GetContentsHandler{
		"foo.txt": []byte("foo.txt"),
		"bar.txt": []byte("bar.txt"),
	}
	assert.Equal(t, expectedGC, c)

	expectedBT := &GetByTypeHandler{
		Files: []string{"foo.txt", "bar.txt"},
		Dirs:  []string{"foo/", "bar/"},
	}
	assert.Equal(t, expectedBT, bt)
}

func setupForHandlersTest() func() {
	restoreStat := Stat
	restoreReadFile := ReadFile
	Stat = mockStat
	ReadFile = mockReadFileAsName
	return func() {
		ReadFile = restoreReadFile
		Stat = restoreStat
	}
}
