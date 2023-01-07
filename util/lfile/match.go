package lfile

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/adamcolton/luce/util/filter"
)

// Match will walk a directory and match files and subdirectories. SkipDir can
// be used to skip directories and their contents. The distinction between
// skipped directories and directories that not found is that when a directory
// is skipped, all of it's contents are skipped, when a directory is not found,
// it's contents are visited, but the directory itself will not be value
// returned by the iterator.
type Match struct {
	SkipDir filter.Filter[string]
	Find    struct {
		File, Dir filter.Filter[string]
	}
}

// NewMatch makes a new instance of Match
func NewMatch(findFile, findDir, skipDir filter.Filter[string]) Match {
	m := Match{
		SkipDir: skipDir,
	}
	m.Find.File = findFile
	m.Find.Dir = findDir
	return m
}

// MustRegexMatch makes an instance of Match using regex for all the values.
func MustRegexMatch(findFile, findDir, skipDir string) Match {
	var m Match
	if findFile != "" {
		m.Find.File = filter.MustRegex(findFile)
	}
	if findDir != "" {
		m.Find.Dir = filter.MustRegex(findDir)
	}
	if skipDir != "" {
		m.SkipDir = filter.MustRegex(skipDir)
	}
	return m
}

// MatchRoot combines a Match instance with a Root. It fulfills IteratorSource.
type MatchRoot struct {
	Match
	Root string
}

// Root to Match against, return MatchRoot which fulfills IteratorSource.
func (m Match) Root(root string) MatchRoot {
	return MatchRoot{
		Match: m,
		Root:  root,
	}
}

type matchRootIter struct {
	MatchRoot
	files []string
	err   error
	done  bool
	data  []byte
	info  os.FileInfo
	idx   int
}

// Iterator fulfills IteratorSource returning an Iterator to iterate over all
// the matches starting from the root.
func (mr MatchRoot) Iterator() (i Iterator, done bool) {
	mri := &matchRootIter{
		MatchRoot: mr,
	}
	return mri, mri.Reset()
}

func (mri *matchRootIter) Idx() int {
	return mri.idx
}

func (mri *matchRootIter) Next() (string, bool) {
	mri.data = nil
	ln := len(mri.files) - 1
	var path string
	updatePath := func() {
		if ln >= 0 {
			path = mri.files[ln]
		}
	}
	updatePath()
	var done bool
	defer func() {
		if !mri.done {
			mri.idx++
		}
	}()
	for mri.done = ln < 0; !mri.done; mri.done = ln < 0 {
		doAppend := true
		if mri.info == nil {
			done, doAppend = mri.checkFilters(path)
			if done {
				break
			}
		}

		mri.files = mri.files[:ln]
		ln--
		if doAppend {
			ln += mri.appendFiles(path)
		}
		mri.info = nil
		updatePath()
		continue
	}
	return mri.Path(), mri.done
}

func (mri *matchRootIter) checkFilters(path string) (done, doAppend bool) {
	mri.info, mri.err = Stat(path)
	if mri.err != nil {
		mri.done = true
		return true, false
	} else if mri.info.IsDir() {
		doAppend = mri.SkipDir == nil || !mri.SkipDir(path)
		return (doAppend && mri.Find.Dir != nil && mri.Find.Dir(path)), doAppend
	}
	return (mri.Find.File != nil && mri.Find.File(path)), false

}

func (mri *matchRootIter) appendFiles(path string) int {
	files, err := readDirNames(path)
	if err != nil {
		mri.err, mri.done = err, true
		return 0
	}
	for _, f := range files {
		mri.files = append(mri.files, filepath.Join(path, f))
	}
	return len(files)
}

func (mri *matchRootIter) Path() string {
	if mri.done {
		return ""
	}
	return mri.path()
}

func (mri *matchRootIter) path() string {
	return mri.files[len(mri.files)-1]
}

func (mri *matchRootIter) Done() bool {
	return mri.done
}

func (mri *matchRootIter) Cur() (path string, done bool) {
	return mri.Path(), mri.done
}

func (mri *matchRootIter) Data() []byte {
	if mri.data == nil && !mri.done {
		mri.data, mri.err = ReadFile(mri.path())
		mri.done = mri.err != nil
	}
	return mri.data
}

func (mri *matchRootIter) Err() error {
	return mri.err
}

func (mri *matchRootIter) Stat() os.FileInfo {
	return mri.info
}

func (mri *matchRootIter) Reset() bool {
	mri.files = mri.files[:0]
	mri.info = nil
	mri.appendFiles(mri.Root)
	mri.done = mri.done || len(mri.files) == 0
	mri.Next()
	mri.idx = 0

	return mri.done
}

var readDirNames = func(dirname string) ([]string, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	sort.Strings(names)
	return names, nil
}
