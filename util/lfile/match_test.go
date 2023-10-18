package lfile_test

import (
	"reflect"
	"testing"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/lfile"
	"github.com/adamcolton/luce/util/lfile/lfilemock"
	"github.com/stretchr/testify/assert"
)

func TestNewMatch(t *testing.T) {
	f := filter.EQ("foo.txt")
	d := filter.EQ("/dir/")
	h := filter.EQ("/.hidden/")
	m := lfile.NewMatch(f, d, h)
	assert.Equal(t, reflect.ValueOf(m.Find.File).Pointer(), reflect.ValueOf(f).Pointer())
	assert.Equal(t, reflect.ValueOf(m.Find.Dir).Pointer(), reflect.ValueOf(d).Pointer())
	assert.Equal(t, reflect.ValueOf(m.SkipDir).Pointer(), reflect.ValueOf(h).Pointer())

	m = lfile.MustRegexMatch(lfile.Exts(false, "txt"), `\/.*dir\/`, `\/\..*/`)
	assert.True(t, m.Find.File("foo.txt"))
	assert.False(t, m.Find.File("/foo.txt/bar"))
	assert.True(t, m.Find.Dir("/mydir/"))
	assert.False(t, m.Find.Dir("/adirs/"))
	assert.True(t, m.SkipDir("/.hidden/"))
	assert.False(t, m.SkipDir("/not.hidden/"))
}

type idxCheck struct {
	shouldBe int
	t        *testing.T
}

func (i *idxCheck) HandleIter(itr lfile.Iterator) {
	assert.Equal(i.t, i.shouldBe, itr.Idx())
	i.shouldBe++
}

func TestMRILMock(t *testing.T) {
	repo := lfilemock.Parse(map[string]any{
		".": []string{"a", "aa", "b", "c"},
		"dir": map[string]any{
			".":    []string{"d", "e", "f"},
			"dir1": []string{"g", "h", "hh", "i"},
			"dir2": []string{"j", "k"},
		},
		".hidden1": []string{"x", "y", "z"},
	}).Repository()

	tt := map[string]struct {
		lfile.Match
		expectedFiles slice.Slice[string]
		expectedDirs  slice.Slice[string]
	}{
		"basic": {
			Match: lfile.MustRegexMatch(`\/[a-z]$`, `[0-9]$`, lfile.FilterHidden),
			expectedFiles: slice.Slice[string]{
				"/dir/f", "/dir/dir2/k", "/dir/dir2/j", "/dir/e", "/dir/dir1/i",
				"/dir/dir1/h", "/dir/dir1/g", "/dir/d", "/c", "/b", "/a",
			},
			expectedDirs: slice.Slice[string]{"/dir/dir2", "/dir/dir1"},
		},
		"no-dirs": {
			Match: lfile.MustRegexMatch(`\/[a-z]$`, ``, lfile.FilterHidden),
			expectedFiles: []string{
				"/dir/f", "/dir/dir2/k", "/dir/dir2/j", "/dir/e", "/dir/dir1/i",
				"/dir/dir1/h", "/dir/dir1/g", "/dir/d", "/c", "/b", "/a",
			},
		},
		"no-files": {
			Match:        lfile.MustRegexMatch(``, `[0-9]$`, lfile.FilterHidden),
			expectedDirs: []string{"/dir/dir2", "/dir/dir1"},
		},
		"no-filter": {
			Match: lfile.MustRegexMatch(`\/[a-z]$`, `[0-9]$`, ``),
			expectedFiles: []string{
				"/.hidden1/z", "/.hidden1/y", "/.hidden1/x", "/dir/f",
				"/dir/dir2/k", "/dir/dir2/j", "/dir/e", "/dir/dir1/i",
				"/dir/dir1/h", "/dir/dir1/g", "/dir/d", "/c", "/b", "/a",
			},
			expectedDirs: []string{"/.hidden1", "/dir/dir2", "/dir/dir1"},
		},
	}

	lt := slice.LT[string]()
	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			r := tc.Match.Root("/")
			r.FSReader = repo
			got := slice.FromIterFactory(r.Factory, nil)
			got.Sort(lt)
			expected := slice.New(append(slice.New(tc.expectedDirs).Clone(), tc.expectedFiles...))
			expected.Sort(lt)
			assert.Equal(t, expected, got)

			i, _ := r.Iterator()
			byType := &lfile.GetByTypeHandler{}
			contents := lfile.GetContentsHandler{}

			idxChecker := &idxCheck{t: t}
			err := lfile.RunHandler(i, lfile.MultiHandler{byType, contents, idxChecker})
			assert.NoError(t, err)
			assert.True(t, i.Done())
			assert.Equal(t, "", i.Path())

			assert.Equal(t, tc.expectedFiles.Sort(lt), slice.New(byType.Files).Sort(lt))
			assert.Equal(t, tc.expectedDirs.Sort(lt), slice.New(byType.Dirs).Sort(lt))
			expectedContents := make(lfile.GetContentsHandler, len(tc.expectedFiles))
			for _, f := range tc.expectedFiles {
				_, n := lfile.Name(f)
				expectedContents[f] = []byte(n)
			}
			assert.Equal(t, expectedContents, contents)
		})
	}
}

func TestMRIErr(t *testing.T) {
	dir := lfilemock.Parse(map[string]any{
		".": []string{"a", "aa", "b", "c"},
		"dir": map[string]any{
			".":    []string{"d", "e", "f"},
			"dir1": []string{"g", "h", "hh", "i"},
			"dir2": []string{"j", "k"},
		},
		".hidden1": []string{"x", "y", "z"},
	})
	dir.Err = lerr.Str("test error")

	mr := lfile.MustRegexMatch(`\/[a-z]$`, `[0-9]$`, lfile.FilterHidden).
		Root("/")
	mr.FSReader = dir.Repository()
	i, _ := mr.Iterator()

	assert.Equal(t, "test error", i.Err().Error())
	_, done := i.Next()
	assert.True(t, done)
}
