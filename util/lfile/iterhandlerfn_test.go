package lfile_test

import (
	"testing"

	"github.com/adamcolton/luce/util/lfile"
	"github.com/stretchr/testify/assert"
)

func TestIterHandlerFn(t *testing.T) {
	var got []string
	var handler lfile.IterHandler = lfile.IterHandlerFn(func(i lfile.Iterator) {
		got = append(got, i.Path())
	})

	// only the paths are used, so no file has to exist
	assert.NoError(t, lfile.RunHandlerSource(lfile.Paths{"a.txt", "dir/b.txt"}, handler))
	assert.Equal(t, []string{"a.txt", "dir/b.txt"}, got)
}
