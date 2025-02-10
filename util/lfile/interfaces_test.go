package lfile_test

import (
	"embed"
	"io/fs"
	"os"
	"testing"
	"testing/fstest"

	"github.com/adamcolton/luce/util/lfile"
	"github.com/stretchr/testify/assert"
)

// The file systems of the standard library fulfill the interfaces, so a caller
// can pass an *os.File, an embed.FS or an fstest.MapFS without an adapter.
var (
	_ lfile.Dir      = (*os.File)(nil)
	_ lfile.File     = (*os.File)(nil)
	_ lfile.CoreFS   = embed.FS{}
	_ lfile.FSReader = fstest.MapFS{}
	_ lfile.FSOpener = os.DirFS(".")
	_ lfile.CoreFS   = os.DirFS(".").(lfile.CoreFS)
)

func TestFSOpenerIsFSFS(t *testing.T) {
	var fsys fs.FS = fstest.MapFS{"a.txt": {Data: []byte("hi")}}
	var opener lfile.FSOpener = fsys
	f, err := opener.Open("a.txt")
	assert.NoError(t, err)
	assert.NoError(t, f.Close())
}

func TestCoreFSOverMapFS(t *testing.T) {
	var cfs lfile.CoreFS = fstest.MapFS{
		"dir/a.txt": {Data: []byte("hi")},
	}

	data, err := cfs.ReadFile("dir/a.txt")
	assert.NoError(t, err)
	assert.Equal(t, "hi", string(data))

	entries, err := cfs.ReadDir("dir")
	assert.NoError(t, err)
	if assert.Len(t, entries, 1) {
		assert.Equal(t, "a.txt", entries[0].Name())
	}
}

func TestFSReaderOverMapFS(t *testing.T) {
	var fsr lfile.FSReader = fstest.MapFS{"a.txt": {Data: []byte("hi")}}
	info, err := fsr.Stat("a.txt")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), info.Size())
}
