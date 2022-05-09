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

	fs := Filenames{"foo.txt", "bar.txt"}
	c := 0
	for i, done := fs.Iter(true); !done; done = i.Next() {
		c++
		assert.Equal(t, fs[i.Index], string(i.Data))
	}
	assert.Equal(t, len(fs), c)
}

func setupForTestIter() func() {
	restoreReadFile := ReadFile
	ReadFile = mockReadFileAsName
	return func() { ReadFile = restoreReadFile }
}
