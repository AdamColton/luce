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
