package lfile_test

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/lfile"
	"github.com/stretchr/testify/assert"
)

func matchedFiles(t *testing.T, mr lfile.MatchRoot) slice.Slice[string] {
	t.Helper()
	files := lfile.GetFiles("", nil)
	assert.NoError(t, lfile.RunHandlerSource(mr, files))
	return files.Matches.Sort(slice.LT[string]())
}

func TestMatchRootPaths(t *testing.T) {
	m := lerr.Must(lfile.RegexMatch(`\.txt$`, "", ""))
	fsys := fstest.MapFS{
		"a.txt":          {},
		"docs/b.txt":     {},
		"docs/sub/c.txt": {},
		"docs/d.md":      {},
	}

	tt := map[string]struct {
		root     []string
		expected slice.Slice[string]
	}{
		// the root of an io/fs file system is ".", so an empty root is
		"empty root":    {[]string{""}, []string{"a.txt", "docs/b.txt", "docs/sub/c.txt"}},
		"no root":       {nil, []string{"a.txt", "docs/b.txt", "docs/sub/c.txt"}},
		"dot":           {[]string{"."}, []string{"a.txt", "docs/b.txt", "docs/sub/c.txt"}},
		"dir":           {[]string{"docs"}, []string{"docs/b.txt", "docs/sub/c.txt"}},
		"trailing":      {[]string{"docs/"}, []string{"docs/b.txt", "docs/sub/c.txt"}},
		"joined":        {[]string{"docs", "sub"}, []string{"docs/sub/c.txt"}},
		"joined, slash": {[]string{"docs/", "/sub/"}, []string{"docs/sub/c.txt"}},
	}
	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			mr := m.Root(tc.root...)
			mr.CoreFS = fsys
			assert.Equal(t, tc.expected, matchedFiles(t, mr))
		})
	}
}

func TestMatchRootOS(t *testing.T) {
	dir := t.TempDir()
	assert.NoError(t, os.Mkdir(filepath.Join(dir, "sub"), 0o755))
	assert.NoError(t, os.WriteFile(filepath.Join(dir, "a.txt"), nil, 0o644))
	assert.NoError(t, os.WriteFile(filepath.Join(dir, "sub", "b.txt"), nil, 0o644))
	m := lerr.Must(lfile.RegexMatch(`\.txt$`, "", ""))
	slash := filepath.ToSlash(dir)

	// the default file system is the OS
	got := matchedFiles(t, m.Root(dir))
	assert.Equal(t, slice.Slice[string]{slash + "/a.txt", slash + "/sub/b.txt"}, got)

	// an io/fs view of a directory takes an empty root
	mr := m.Root("")
	mr.CoreFS = os.DirFS(dir).(lfile.CoreFS)
	got = matchedFiles(t, mr)
	assert.Equal(t, slice.Slice[string]{"a.txt", "sub/b.txt"}, got)
}
