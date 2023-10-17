package lfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetByTypeHandler(t *testing.T) {
	fs := Paths{"foo/", "foo.txt", "bar.txt", "bar/"}.FS(nameFS{})
	bt := &GetByTypeHandler{}
	err := RunHandlerSource(fs, bt)
	assert.NoError(t, err)

	expected := &GetByTypeHandler{
		Files: []string{"foo.txt", "bar.txt"},
		Dirs:  []string{"foo/", "bar/"},
	}
	assert.Equal(t, expected, bt)
}
