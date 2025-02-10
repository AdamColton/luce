package lfile

import (
	"bytes"
	"testing"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

// treePaths are a small tree for nameFS, in which a name that ends in a slash is a
// directory.
var treePaths = Paths{"root/a/", "root/a/b.txt", "root/a/c/", "root/a/c/d.txt", "root/e.txt", "root/f/"}

func TestGetType(t *testing.T) {
	src := treePaths.FS(nameFS{})

	files := GetFiles("root/", nil)
	assert.NoError(t, RunHandlerSource(src, files))
	assert.Equal(t, slice.Slice[string]{"a/b.txt", "a/c/d.txt", "e.txt"}, files.Matches)

	dirs := GetDirs("root/", nil)
	assert.NoError(t, RunHandlerSource(src, dirs))
	assert.Equal(t, slice.Slice[string]{"a/", "a/c/", "f/"}, dirs.Matches)

	// without a prefix the paths are whole, and a buffer is used
	buf := make([]string, 0, 10)
	dirs = GetDirs("", buf)
	assert.NoError(t, RunHandlerSource(src, dirs))
	assert.Equal(t, "root/a/", dirs.Matches[0])
	assert.Same(t, &buf[:1][0], &dirs.Matches[0])

	// a prefix that is not there is not removed
	files = GetFiles("other/", nil)
	assert.NoError(t, RunHandlerSource(src, files))
	assert.Equal(t, "root/a/b.txt", files.Matches[0])
}

func TestFilesTree(t *testing.T) {
	ft := NewFilesTree("root/")
	assert.NoError(t, RunHandlerSource(treePaths.FS(nameFS{}), ft))

	root := ft.Root()
	assert.True(t, root.IsDir())
	assert.Equal(t, "", root.Name())
	assert.Len(t, root.Children(), 3)

	a := root.Children()["a"]
	assert.True(t, a.IsDir())
	assert.Equal(t, "a", a.Name())
	assert.False(t, a.Children()["b.txt"].IsDir())
	assert.Nil(t, a.Children()["b.txt"].Children())
	assert.True(t, a.Children()["c"].IsDir())
	assert.False(t, root.Children()["e.txt"].IsDir())
	assert.True(t, root.Children()["f"].IsDir(), "an empty directory is a directory")

	// written in order
	buf := bytes.NewBuffer(nil)
	root.Write(buf, "", "  ")
	assert.Equal(t, " a\n   b.txt\n   c\n     d.txt\n e.txt\n f\n", buf.String())

	// the prefix is the path that is cut, so it can be the whole path
	ft = NewFilesTree("")
	assert.NoError(t, RunHandlerSource(Paths{"x/", "x/y.txt"}.FS(nameFS{}), ft))
	assert.Equal(t, "y.txt", ft.Root().Children()["x"].Children()["y.txt"].Name())
}

func TestFileTreeNext(t *testing.T) {
	root := NewFilesTree("").Root()

	// a nil context creates files
	n, found := root.Next("a", true, nil)
	assert.True(t, found)
	assert.False(t, n.IsDir())
	_, found = root.Next("b", false, nil)
	assert.False(t, found)
	again, found := root.Next("a", false, nil)
	assert.True(t, found)
	assert.Same(t, n, again)

	// a context that says how far along it is makes directories on the way
	ctx := &FileTreeCtx{depth: 2}
	n, _ = root.Next("d", true, ctx)
	assert.True(t, n.IsDir())
	assert.Equal(t, 1, ctx.depth)
	n, _ = n.Next("f", true, ctx)
	assert.False(t, n.IsDir())
	ctx = &FileTreeCtx{depth: 1, isDir: true}
	n, _ = root.Next("g", true, ctx)
	assert.True(t, n.IsDir())
}
