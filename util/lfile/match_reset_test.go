package lfile_test

import (
	"testing"
	"testing/fstest"

	"github.com/adamcolton/luce/util/lfile"
	"github.com/stretchr/testify/assert"
)

// txtRoot matches .txt files below root of cfs.
func txtRoot(cfs lfile.CoreFS, root string) lfile.MatchRoot {
	m, err := lfile.RegexMatch(`\.txt$`, "", "")
	if err != nil {
		panic(err)
	}
	mr := m.Root(root)
	mr.CoreFS = cfs
	return mr
}

func collect(i lfile.Iterator, done bool) []string {
	var out []string
	for ; !done; _, done = i.Next() {
		out = append(out, i.Path())
	}
	return out
}

func TestMatchRootResetAfterEnd(t *testing.T) {
	m := fstest.MapFS{
		"a.txt": {Data: []byte("a")},
		"b.txt": {Data: []byte("b")},
	}
	i, done := txtRoot(m, ".").Iterator()
	first := collect(i, done)
	assert.ElementsMatch(t, []string{"a.txt", "b.txt"}, first)
	assert.True(t, i.Done())

	// Reset after the end starts again
	assert.False(t, i.Reset())
	assert.Equal(t, 0, i.Idx())
	assert.Equal(t, first, collect(i, false))
}

func TestMatchRootResetAfterError(t *testing.T) {
	m := fstest.MapFS{
		"a.txt": {Data: []byte("a")},
		"b.txt": {Data: []byte("b")},
	}
	i, done := txtRoot(m, ".").Iterator()
	assert.False(t, done)

	// remove a file that was already listed so it can't be read
	delete(m, i.Path())
	assert.Nil(t, i.Data())
	assert.Error(t, i.Err())
	assert.True(t, i.Done())

	assert.False(t, i.Reset())
	assert.NoError(t, i.Err())
	assert.Len(t, collect(i, false), 1)

	// A root that can't be read is an error until it exists
	m2 := fstest.MapFS{}
	i, done = txtRoot(m2, "dir").Iterator()
	assert.True(t, done)
	assert.Error(t, i.Err())
	m2["dir/c.txt"] = &fstest.MapFile{}
	assert.False(t, i.Reset())
	assert.NoError(t, i.Err())
	assert.Equal(t, "dir/c.txt", i.Path())
}
