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
