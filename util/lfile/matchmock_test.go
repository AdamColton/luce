package lfile_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/util/lfile"
	"github.com/adamcolton/luce/util/lfile/lfilemock"
	"github.com/stretchr/testify/assert"
)

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

	f, err := repo.Open("/dir/f")
	assert.NoError(t, err)
	assert.NotNil(t, f)

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
		// "no-filter": {
		// 	Match: lfile.MustRegexMatch(`\/[a-z]$`, `[0-9]$`, ``),
		// 	expectedFiles: []string{
		// 		"/.hidden1/z", "/.hidden1/y", "/.hidden1/x", "/dir/f",
		// 		"/dir/dir2/k", "/dir/dir2/j", "/dir/e", "/dir/dir1/i",
		// 		"/dir/dir1/h", "/dir/dir1/g", "/dir/d", "/c", "/b", "/a",
		// 	},
		// 	expectedDirs: []string{"/.hidden1", "/dir/dir2", "/dir/dir1"},
		// },
	}

	lt := slice.LT[string]()

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			r := tc.Match.
				Root("/")
			r.Repository = repo
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
