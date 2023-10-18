package lfile

import (
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountAll(t *testing.T) {
	tt := map[string]struct {
		s                string
		w                []*regexp.Regexp
		expectedfoundAll bool
		expectedMatches  [][][]int
		expectedCounts   []int
		expectedSum      int
	}{
		"basic": {
			s:                "foo bar test baz foo",
			w:                RegexWords([]string{"Foo", "Bar", "Baz"}),
			expectedfoundAll: true,
			expectedCounts:   []int{2, 1, 1},
			expectedMatches: [][][]int{
				{{0, 3}, {17, 20}},
				{{4, 7}},
				{{13, 16}},
			},
			expectedSum: 4,
		},
		"miss-one": {
			s:                "foo bar test foo",
			w:                RegexWords([]string{"Foo", "Bar", "Baz"}),
			expectedfoundAll: false,
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			f, m, c, s := CountAll(tc.s, tc.w)
			assert.Equal(t, tc.expectedfoundAll, f)
			assert.Equal(t, tc.expectedMatches, m)
			assert.Equal(t, tc.expectedCounts, c)
			assert.Equal(t, tc.expectedSum, s)
		})
	}
}

func TestIterSearch(t *testing.T) {
	files, restore := setupForTestIterSearch()
	defer restore()

	search := &RegexSearch{
		Terms: RegexWords([]string{"test", "foo"}),
	}
	err := RunHandlerSource(files, search)
	assert.NoError(t, err)

	if assert.Len(t, search.Results, 1) {
		r := search.Results[0]
		assert.Equal(t, "foo.txt", r.File)
		assert.Equal(t, 2, r.Sum)
		assert.Equal(t, []int{1, 1}, r.Counts)
		assert.Equal(t, [][][]int{{{19, 23}}, {{8, 11}}}, r.Matches)
	}
}

func setupForTestIterSearch() (Paths, func()) {
	mocks := map[string]struct {
		info    mockFileInfo
		data    string
		skipDir bool
	}{
		"foo.txt": {
			info: mockFileInfo{
				isDir: false,
			},
			data: "this is foo.txt, a test",
		},
		"bar.txt": {
			info: mockFileInfo{
				isDir: false,
			},
			data: "contains no search words",
		},
		"baz.txt": {
			info: mockFileInfo{
				isDir: false,
			},
			data: "this test contains one search word",
		},
		"/mydir/": {
			info: mockFileInfo{
				isDir: true,
			},
		},
		"/.hidden/": {
			info: mockFileInfo{
				isDir: true,
			},
			skipDir: true,
		},
	}

	restoreReadFile := ReadFile
	restore := func() {
		ReadFile = restoreReadFile
	}

	ReadFile = func(filename string) ([]byte, error) {
		return []byte(mocks[filename].data), nil
	}
	Stat = func(filename string) (os.FileInfo, error) {
		return mocks[filename].info, nil
	}
	files := make(Paths, 0, len(mocks))
	for f := range mocks {
		files = append(files, f)
	}

	return files, restore
}
