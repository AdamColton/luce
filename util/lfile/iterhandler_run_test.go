package lfile_test

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/adamcolton/luce/util/lfile"
	"github.com/stretchr/testify/assert"
)

type funcHandler func(lfile.Iterator)

func (f funcHandler) HandleIter(i lfile.Iterator) { f(i) }

// TestRunHandlerErrors shows that the error that ends an iteration is what the
// run returns, and that the handler is not called for the path that failed.
func TestRunHandlerErrors(t *testing.T) {
	files := fstest.MapFS{"a.txt": {Data: []byte("a")}, "c.txt": {Data: []byte("c")}}
	var read []string
	handler := funcHandler(func(i lfile.Iterator) {
		if len(i.Data()) > 0 {
			read = append(read, i.Path())
		}
	})

	src := lfile.Paths{"a.txt", "missing.txt", "c.txt"}.FS(files)
	err := lfile.RunHandlerSource(src, handler)
	assert.ErrorIs(t, err, fs.ErrNotExist)
	assert.Equal(t, []string{"a.txt"}, read, "the iteration ends at the file that fails")

	// RunHandler starts again from the beginning
	i, _ := src.Iterator()
	read = nil
	assert.ErrorIs(t, lfile.RunHandler(i, handler), fs.ErrNotExist)
	assert.Equal(t, []string{"a.txt"}, read)
}
