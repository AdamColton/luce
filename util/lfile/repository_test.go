package lfile_test

import (
	"testing"

	"github.com/adamcolton/luce/util/lfile"
	"github.com/adamcolton/luce/util/lfile/lfilemock"
	"github.com/stretchr/testify/assert"
)

func TestRepositorySize(t *testing.T) {
	r := lfilemock.Parse(map[string]any{
		"file1.txt": "this is test file 1",
		"dir1": map[string]any{
			"file2.bin": []byte{3, 1, 4, 1, 5, 9, 2, 6, 5},
		},
		"dir2": map[string]any{
			"dir3":      map[string]any{},
			"file4.txt": "this is file 4",
		},
	}).Repository()

	size, err := lfile.Size(r, "/")
	assert.NoError(t, err)
	assert.Equal(t, int64(42), size)
}
