package lfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockReadFileAsName(filename string) ([]byte, error) {
	return []byte(filename), nil
}

func TestIter(t *testing.T) {
	restore := setupForTestIter()
	defer restore()

	fs := Paths{"foo.txt", "bar.txt"}
	c := 0
	i, done := fs.Iterator()
	for ; !done; done = i.Next() {
		c++
		assert.Equal(t, fs[i.(*pathsIterator).Index], string(i.Data()))
		assert.False(t, i.Done())
	}
	assert.True(t, i.Done())
	assert.Equal(t, len(fs), c)
}

func setupForTestIter() func() {
	restoreReadFile := ReadFile
	ReadFile = mockReadFileAsName
	return func() { ReadFile = restoreReadFile }
}
