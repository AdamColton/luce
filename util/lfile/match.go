package lfile

import (
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/liter"
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
		m.Find.File = lerr.Must(filter.Regex(findFile))
	}
	if findDir != "" {
		m.Find.Dir = lerr.Must(filter.Regex(findDir))
	}
	if skipDir != "" {
		m.SkipDir = lerr.Must(filter.Regex(skipDir))
	}
	return m
}

// MatchRoot combines a Match instance with a Root. It fulfills IteratorSource.
// Repository allows for testing. If Repository is nil, OSRepository is used.
type MatchRoot struct {
	Match
	Root string
	Repository
}

// Root to Match against, return MatchRoot which fulfills IteratorSource.
func (m Match) Root(root string) MatchRoot {
	return MatchRoot{
		Match:      m,
		Root:       root,
		Repository: OSRepository{},
	}
}

type matchRootIter struct {
	MatchRoot
	path  string
	files slice.Slice[string]
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

func (mr MatchRoot) Factory() (i liter.Iter[string], str string, done bool) {
	mri := &matchRootIter{
		MatchRoot: mr,
	}
	done = mri.Reset()
	if !done {
		str = mri.path
	}
	i = mri
	return
}

func (mri *matchRootIter) Idx() int {
	return mri.idx
}

func (mri *matchRootIter) Next() (string, bool) {
	mri.data = nil
	done := func() bool {
		mri.done = mri.done || len(mri.files) == 0
		return mri.done
	}
	for !done() {
		mri.path, mri.files = mri.files.Pop()
		mri.info, mri.err = mri.Repository.Stat(mri.path)
		foundNext, doAppend := mri.checkFilters()
		if doAppend {
			mri.appendFiles()
		}
		if foundNext {
			break
		}
	}
	if !mri.done {
		mri.idx++
	}
	return mri.Path(), mri.done
}

func (mri *matchRootIter) checkFilters() (foundNext, doAppend bool) {
	if mri.err != nil {
		mri.done = true
		return true, false
	} else if mri.info.IsDir() {
		doAppend = mri.SkipDir == nil || !mri.SkipDir(mri.path)
		foundNext = (doAppend && mri.Find.Dir != nil && mri.Find.Dir(mri.path))
		return
	}
	doAppend = false
	foundNext = (mri.Find.File != nil && mri.Find.File(mri.path))
	return
}

func (mri *matchRootIter) appendFiles() {
	files, err := readDirNames(mri.Repository, mri.path)
	if err != nil {
		mri.err, mri.done = err, true
		return
	}
	for _, f := range files {
		mri.files = append(mri.files, filepath.Join(mri.path, f))
	}
}

func (mri *matchRootIter) Path() string {
	if mri.done {
		return ""
	}
	return mri.path
}

func (mri *matchRootIter) Done() bool {
	return mri.done
}

func (mri *matchRootIter) Cur() (path string, done bool) {
	return mri.Path(), mri.done
}

func (mri *matchRootIter) Data() []byte {
	if mri.data == nil && !mri.done {
		mri.data, mri.err = mri.Repository.ReadFile(mri.path)
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
	mri.path = mri.Root
	mri.appendFiles()
	mri.done = mri.done || len(mri.files) == 0
	mri.Next()
	mri.idx = 0

	return mri.done
}

var readDirNames = func(r Repository, dirname string) ([]string, error) {
	f, err := r.Open(dirname)
	if err != nil {
		return nil, err
	}
	names, err := f.Readdirnames(-1)
	f.Close()

	if err != nil && !lerr.Except(err, io.EOF) {
		return nil, err
	}
	sort.Strings(names)
	return names, nil
}
