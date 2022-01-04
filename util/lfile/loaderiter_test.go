package lfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIter(t *testing.T) {
	restore := ReadFile
	defer func() { ReadFile = restore }()
	ReadFile = func(filename string) ([]byte, error) {
		return []byte(filename), nil
	}

	fs := Filenames{"foo.txt", "bar.txt"}
	c := 0
	for i, done := fs.Iter(true); !done; done = i.Next() {
		c++
		assert.Equal(t, fs[i.Index], string(i.Data))
	}
	assert.Equal(t, len(fs), c)
}
