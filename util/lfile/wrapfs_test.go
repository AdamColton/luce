package lfile_test

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/adamcolton/luce/util/lfile"
	"github.com/adamcolton/luce/util/upgrade"
	"github.com/stretchr/testify/assert"
)

// openOnly is an fs.FS with nothing but Open, like an archive/zip reader.
type openOnly struct {
	m fstest.MapFS
}

func (o openOnly) Open(name string) (fs.File, error) { return o.m.Open(name) }

func TestWrapFS(t *testing.T) {
	zipLike := openOnly{fstest.MapFS{
		"docs/a.txt": {Data: []byte("hello")},
		"docs/b.txt": {Data: []byte("world")},
	}}
	_, hasStat := any(zipLike).(lfile.FSStater)
	assert.False(t, hasStat)

	fsr := lfile.WrapFS(zipLike)

	data, err := fsr.ReadFile("docs/a.txt")
	assert.NoError(t, err)
	assert.Equal(t, "hello", string(data))

	info, err := fsr.Stat("docs/b.txt")
	assert.NoError(t, err)
	assert.Equal(t, int64(5), info.Size())

	entries, err := fsr.ReadDir("docs")
	assert.NoError(t, err)
	assert.Len(t, entries, 2)

	// The files it opens are Dirs.
	f, err := fsr.Open("docs")
	assert.NoError(t, err)
	defer f.Close()
	dir, isDir := f.(lfile.Dir)
	assert.True(t, isDir)
	names, err := dir.Readdirnames(0)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"a.txt", "b.txt"}, names)

	_, err = fsr.ReadFile("missing")
	assert.ErrorIs(t, err, fs.ErrNotExist)

	w, ok := upgrade.To[openOnly](fsr)
	assert.True(t, ok)
	assert.Equal(t, zipLike, w)
}

func TestWrapFSUsesFileSystemMethods(t *testing.T) {
	m := fstest.MapFS{"a.txt": {Data: []byte("hi")}}

	// A MapFS is already an FSReader and comes back untouched.
	assert.Equal(t, lfile.FSReader(m), lfile.WrapFS(m))

	// A CoreFS with no Stat gets the fallbacks of fs.FS.
	core := newCoreOnly(m)
	fsr := lfile.WrapFS(core)
	data, err := fsr.ReadFile("a.txt")
	assert.NoError(t, err)
	assert.Equal(t, "hi", string(data))
	info, err := fsr.Stat("a.txt")
	assert.NoError(t, err)
	assert.Equal(t, "a.txt", info.Name())
}
