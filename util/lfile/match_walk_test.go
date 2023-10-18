package lfile_test

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/lfile"
	"github.com/stretchr/testify/assert"
)

func TestRegexMatchErrors(t *testing.T) {
	// each expression is compiled, and an empty one is left nil
	m, err := lfile.RegexMatch("", "", "")
	assert.NoError(t, err)
	assert.Nil(t, m.Find.File)
	assert.Nil(t, m.Find.Dir)
	assert.Nil(t, m.SkipDir)

	m, err = lfile.RegexMatch(`\.go$`, `^src`, `^\.`)
	assert.NoError(t, err)
	assert.True(t, m.Find.File("a.go"))
	assert.True(t, m.Find.Dir("src/a"))
	assert.True(t, m.SkipDir(".git"))

	for n, args := range map[string][3]string{
		"file": {"(", "", ""},
		"dir":  {"", "(", ""},
		"skip": {"", "", "("},
	} {
		t.Run(n, func(t *testing.T) {
			_, err := lfile.RegexMatch(args[0], args[1], args[2])
			assert.Error(t, err)
		})
	}
}

func walk(t *testing.T, m lfile.Match, cfs lfile.CoreFS, root string) []string {
	t.Helper()
	mr := m.Root(root)
	mr.CoreFS = cfs
	var got []string
	i, done := mr.Iterator()
	for ; !done; _, done = i.Next() {
		got = append(got, i.Path())
	}
	assert.NoError(t, i.Err())
	return got
}

func TestMatchWalk(t *testing.T) {
	tree := fstest.MapFS{
		"a.txt":        {},
		"b/c.txt":      {},
		"b/d.txt":      {},
		"b/deep/e.txt": {},
		"skip/f.txt":   {},
		"e.txt":        {},
	}
	files := lerr.Must(lfile.RegexMatch(`\.txt$`, "", `^skip$`))

	// the entries of a directory in reverse order of name, and what is in a
	// directory before the entries that come before it
	assert.Equal(t, []string{"e.txt", "b/deep/e.txt", "b/d.txt", "b/c.txt", "a.txt"}, walk(t, files, tree, "."))
	assert.Equal(t, []string{"b/deep/e.txt", "b/d.txt", "b/c.txt"}, walk(t, files, tree, "b"))

	// directories that are found, of which the root is not one
	dirs := lerr.Must(lfile.RegexMatch("", `.*`, ""))
	assert.Equal(t, []string{"skip", "b", "b/deep"}, walk(t, dirs, tree, "."))
	assert.Equal(t, []string{"b/deep"}, walk(t, dirs, tree, "b"))

	// a directory that is skipped is not returned, and neither is what is in it,
	// though a directory that is only not found is entered
	dirsSkip := lerr.Must(lfile.RegexMatch(`\.txt$`, `^b$`, `^skip$`))
	assert.Equal(t, []string{"e.txt", "b", "b/deep/e.txt", "b/d.txt", "b/c.txt", "a.txt"}, walk(t, dirsSkip, tree, "."))
	all := lerr.Must(lfile.RegexMatch("", `.*`, `^b$`))
	assert.Equal(t, []string{"skip"}, walk(t, all, tree, "."))

	// nothing is found without filters
	assert.Empty(t, walk(t, lfile.Match{}, tree, "."))
}

// TestMatchWalkErrors shows that a path that can't be described ends the
// iteration, and that the error is kept.
func TestMatchWalkErrors(t *testing.T) {
	tree := listedFS{fstest.MapFS{"ghosts/a.txt": {}}}
	m := lerr.Must(lfile.RegexMatch(`\.txt$`, "", ""))
	mr := m.Root("ghosts")
	mr.CoreFS = tree

	// listedFS lists "a-ghost" first, which can't be opened; the files are
	// visited in reverse order so a.txt is found before it
	i, done := mr.Iterator()
	assert.False(t, done)
	assert.Equal(t, "ghosts/a.txt", i.Path())
	path, done := i.Next()
	assert.True(t, done)
	assert.Equal(t, "", path)
	assert.ErrorIs(t, i.Err(), fs.ErrNotExist)
	assert.True(t, i.Done())

	// the factory gives the paths
	got := slice.FromIterFactory(mr.Factory, nil)
	assert.Equal(t, slice.Slice[string]{"ghosts/a.txt"}, got)
}
