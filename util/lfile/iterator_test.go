package lfile_test

import (
	"testing"
	"testing/fstest"

	"github.com/adamcolton/luce/util/lfile"
	"github.com/stretchr/testify/assert"
)

type deletedIter struct {
	lfile.Iterator
}

func (deletedIter) Deleted() bool { return true }

func TestDeleted(t *testing.T) {
	fsys := fstest.MapFS{"notes/a.txt": {Data: []byte("a")}}
	i, _ := lfile.MatchExt(true, "txt").Root("notes").SetCoreFS(fsys).Iterator()
	assert.False(t, lfile.Deleted(i), "an Iterator that is not a Deleter")
	assert.True(t, lfile.Deleted(deletedIter{i}))
}
