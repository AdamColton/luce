package lfile_test

import (
	"testing"
	"testing/fstest"

	"github.com/adamcolton/luce/util/lfile"
	"github.com/stretchr/testify/assert"
)

func TestPathsResetAfterError(t *testing.T) {
	m := fstest.MapFS{"b.txt": {Data: []byte("b")}}
	i, done := lfile.Paths{"a.txt", "b.txt"}.FS(m).Iterator()
	assert.False(t, done)

	// a.txt does not exist, so reading it ends the iteration
	assert.Nil(t, i.Data())
	assert.Error(t, i.Err())
	assert.True(t, i.Done())

	// Reset starts over and forgets the error
	assert.False(t, i.Reset())
	assert.NoError(t, i.Err())
	assert.Equal(t, "a.txt", i.Path())
	path, done := i.Next()
	assert.False(t, done)
	assert.Equal(t, "b.txt", path)
	assert.Equal(t, "b", string(i.Data()))
}
