package lfile

import (
	"os"
	"reflect"
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/filter"
	"github.com/stretchr/testify/assert"
)

func TestNewMatch(t *testing.T) {
	f := filter.EQ("foo.txt")
	d := filter.EQ("/dir/")
	h := filter.EQ("/.hidden/")
	m := NewMatch(f, d, h)
	assert.Equal(t, reflect.ValueOf(m.Find.File).Pointer(), reflect.ValueOf(f).Pointer())
	assert.Equal(t, reflect.ValueOf(m.Find.Dir).Pointer(), reflect.ValueOf(d).Pointer())
	assert.Equal(t, reflect.ValueOf(m.SkipDir).Pointer(), reflect.ValueOf(h).Pointer())

	m = MustRegexMatch(Exts(false, "txt"), `\/.*dir\/`, `\/\..*/`)
	assert.True(t, m.Find.File("foo.txt"))
	assert.False(t, m.Find.File("/foo.txt/bar"))
	assert.True(t, m.Find.Dir("/mydir/"))
	assert.False(t, m.Find.Dir("/adirs/"))
	assert.True(t, m.SkipDir("/.hidden/"))
	assert.False(t, m.SkipDir("/not.hidden/"))
}

// TODO:tests where match has nil values

func TestMRI(t *testing.T) {
	restore := setupForTestMRI()
	defer restore()

	tt := map[string]struct {
		Match
		expectedFiles []string
		expectedDirs  []string
	}{
		"basic": {
			Match: MustRegexMatch(`\/[a-z]$`, `[0-9]$`, FilterHidden),
			expectedFiles: []string{
				"/dir/f", "/dir/dir2/k", "/dir/dir2/j", "/dir/e", "/dir/dir1/i",
				"/dir/dir1/h", "/dir/dir1/g", "/dir/d", "/c", "/b", "/a",
			},
			expectedDirs: []string{"/dir/dir2", "/dir/dir1"},
		},
		"no-dirs": {
			Match: MustRegexMatch(`\/[a-z]$`, ``, FilterHidden),
			expectedFiles: []string{
				"/dir/f", "/dir/dir2/k", "/dir/dir2/j", "/dir/e", "/dir/dir1/i",
				"/dir/dir1/h", "/dir/dir1/g", "/dir/d", "/c", "/b", "/a",
			},
		},
		"no-files": {
			Match:        MustRegexMatch(``, `[0-9]$`, FilterHidden),
			expectedDirs: []string{"/dir/dir2", "/dir/dir1"},
		},
		"no-filter": {
			Match: MustRegexMatch(`\/[a-z]$`, `[0-9]$`, ``),
			expectedFiles: []string{
				"/.hidden1/z", "/.hidden1/y", "/.hidden1/x", "/dir/f",
				"/dir/dir2/k", "/dir/dir2/j", "/dir/e", "/dir/dir1/i",
				"/dir/dir1/h", "/dir/dir1/g", "/dir/d", "/c", "/b", "/a",
			},
			expectedDirs: []string{"/.hidden1", "/dir/dir2", "/dir/dir1"},
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			i, _ := tc.Match.
				Root("/").
				Iterator()
			byType := &GetByTypeHandler{}
			contents := GetContentsHandler{}
			err := RunHandler(i, MultiHandler{byType, contents})
			assert.NoError(t, err)
			assert.True(t, i.Done())
			assert.Equal(t, "", i.Path())
			assert.Equal(t, tc.expectedFiles, byType.Files)
			assert.Equal(t, tc.expectedDirs, byType.Dirs)
			expectedContents := make(GetContentsHandler, len(tc.expectedFiles))
			for _, f := range tc.expectedFiles {
				expectedContents[f] = []byte(f)
			}
			assert.Equal(t, expectedContents, contents)
		})
	}
}

func setupForTestMRI() func() {
	restoreReadDirNames := readDirNames
	restoreStat := Stat
	restoreReadFile := ReadFile

	mockDir := map[string][]string{
		"/":         {"a", "aa", "b", "c", "dir", ".hidden1"},
		"/dir":      {"d", "dir1", "e", "dir2", "f"},
		"/dir/dir1": {"g", "h", "hh", "i"},
		"/dir/dir2": {"j", "k"},
		"/.hidden1": {"x", "y", "z"},
	}
	readDirNames = func(dirname string) ([]string, error) {
		return mockDir[dirname], nil
	}
	Stat = func(name string) (os.FileInfo, error) {
		_, isDir := mockDir[name]
		return mockFileInfo{
			isDir: isDir,
		}, nil
	}
	ReadFile = mockReadFileAsName

	return func() {
		readDirNames, Stat, ReadFile = restoreReadDirNames, restoreStat, restoreReadFile
	}
}

func TestReadDirNames(t *testing.T) {
	names, err := readDirNames(".")
	assert.NoError(t, err)
	expected := []string{
		"dir.go", "dir_test.go", "doc.go", "exts.go", "exts_test.go",
		"handlers.go", "handlers_test.go", "iterator.go", "match.go",
		"match_test.go", "multiglob.go", "multiglob_test.go", "path.go",
		"path_test.go", "paths.go", "paths_test.go", "search.go",
		"search_test.go",
	}
	assert.Equal(t, expected, names)
}

func TestMRIErr(t *testing.T) {
	restoreReadDirNames := readDirNames
	defer func() {
		readDirNames = restoreReadDirNames
	}()

	readDirNames = func(dirname string) ([]string, error) {
		return nil, lerr.Str("test error")
	}

	i, _ := MustRegexMatch(`\/[a-z]$`, `[0-9]$`, FilterHidden).
		Root("/").
		Iterator()

	assert.Equal(t, "test error", i.Err().Error())
	assert.True(t, i.Next())
}
