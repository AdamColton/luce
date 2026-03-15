package lfile_test

import (
	"sort"
	"testing"
	"testing/fstest"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/lfile"
	"github.com/stretchr/testify/assert"
)

func TestSetCoreFS(t *testing.T) {
	m := lerr.Must(lfile.RegexMatch(`\.html$`, "", ""))
	root := m.Root("testdata/templates")

	// an embed.FS can't describe a file, so it is wrapped to be able to
	_, hasStat := any(testdataFS).(lfile.FSStater)
	assert.False(t, hasStat)
	mr := root.SetCoreFS(testdataFS)
	fsr, ok := mr.CoreFS.(lfile.FSReader)
	if assert.True(t, ok) {
		info, err := fsr.Stat("testdata/templates/index.html")
		assert.NoError(t, err)
		assert.Equal(t, "index.html", info.Name())
		f, err := fsr.Open("testdata/templates")
		assert.NoError(t, err)
		_, ok = f.(lfile.Dir)
		assert.True(t, ok, "the files it opens are Dirs")
	}
	w, ok := mr.CoreFS.(interface{ Wrapped() any })
	assert.True(t, ok)
	assert.Equal(t, testdataFS, w.Wrapped())

	// one that can is used as it is
	assert.Equal(t, lfile.OSRepository{}, root.SetCoreFS(lfile.OSRepository{}).CoreFS)
	files := fstest.MapFS{"a.html": {}}
	assert.Equal(t, files, root.SetCoreFS(files).CoreFS)

	// and so is one that has already been wrapped
	assert.Equal(t, mr.CoreFS, mr.SetCoreFS(mr.CoreFS).CoreFS)

	// a minimal CoreFS is wrapped
	core := newCoreOnly(files)
	_, ok = root.SetCoreFS(core).CoreFS.(lfile.FSReader)
	assert.True(t, ok)

	// nil is the default, the OS
	assert.Equal(t, lfile.OSRepository{}, root.SetCoreFS(nil).CoreFS)

	// the MatchRoot it was called on is not changed
	assert.Equal(t, lfile.OSRepository{}, root.CoreFS)

	// the walk gives the same files with the wrapped file system
	got := lfile.GetFiles("", nil)
	assert.NoError(t, lfile.RunHandlerSource(mr, got))
	sort.Strings(got.Matches)
	assert.Equal(t, []string{"testdata/templates/index.html", "testdata/templates/layout/base.html"}, []string(got.Matches))
}
