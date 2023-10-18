package lfile

import (
	"os"
	"path"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/liter"
)

// Match describes which files and directories a walk returns and which
// directories it goes into. Each is a filter on the path, and can be nil: a nil
// Find.File or Find.Dir finds nothing and a nil SkipDir skips nothing. The
// distinction between skipped directories and directories that are not found is
// that when a directory is skipped, none of its contents are visited, while when
// a directory is not found its contents are visited, but the directory itself is
// not returned. A directory that is skipped is never returned, even if Find.Dir
// matches it.
type Match struct {
	// SkipDir is tested against the path of each directory.
	SkipDir filter.Filter[string]
	Find    struct {
		// File is tested against the path of each file and Dir against the path
		// of each directory that is not skipped.
		File, Dir filter.Filter[string]
	}
}

// NewMatch makes a new instance of Match. Any of the filters can be nil.
func NewMatch(findFile, findDir, skipDir filter.Filter[string]) Match {
	m := Match{
		SkipDir: skipDir,
	}
	m.Find.File = findFile
	m.Find.Dir = findDir
	return m
}

// RegexMatch makes an instance of Match using regular expressions for all the
// values, which are matched against the path. An empty expression leaves that
// filter nil. If an expression does not compile the error is returned, along
// with the filters that were made before it.
func RegexMatch(findFile, findDir, skipDir string) (Match, error) {
	var m Match
	var err error
	if findFile != "" {
		m.Find.File, err = filter.Regex(findFile)
		if err != nil {
			return m, err
		}
	}
	if findDir != "" {
		m.Find.Dir, err = filter.Regex(findDir)
		if err != nil {
			return m, err
		}
	}
	if skipDir != "" {
		m.SkipDir, err = filter.Regex(skipDir)
		if err != nil {
			return m, err
		}
	}
	return m, nil
}

// MatchRoot is a Match, the directory to walk and the file system it is in. It
// fulfills IteratorSource. Match.Root creates one that walks the operating
// system, and CoreFS can be changed to walk another file system, such as an
// embed.FS or a mock; it must not be nil.
//
// The paths it returns are slash separated: the Root joined with the names of
// what is found, which are also the paths the filters are tested against. The
// Root itself is not tested or returned. The entries of a directory are visited
// in reverse order of name, and everything in a directory before the entries
// that come before it.
type MatchRoot struct {
	Match
	// Root is the directory to walk. It can be changed, and an empty Root is the
	// current directory.
	Root   string
	CoreFS CoreFS
}

// Root to Match against, return MatchRoot which fulfills IteratorSource. The
// elements are joined with path.Join. The paths that are matched are slash
// separated and are made by joining the root and the name of each entry, so an
// empty root, which is the current directory of the OSRepository and the root
// (".") of an io/fs file system, gives paths with no prefix.
func (m Match) Root(root ...string) MatchRoot {
	return MatchRoot{
		Match:  m,
		Root:   path.Join(root...),
		CoreFS: OSRepository{},
	}
}

// == projects.Code.luce.lfile ==
// [ ] rebuild MatchRoot on fs.WalkDir
//  MatchRoot is its own walker: Match.SkipDir and Match.Find, Iterator.Data and
//  Iterator.Stat for the current path, Reset, and Factory. fs.SkipDir would do
//  for SkipDir and the DirEntry has the Stat, but WalkDir has no Data, Reset or
//  Factory, and does not offer a directory's entries in this order. Check what
//  the later code that uses lfile needs before changing it.

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
// the matches starting from the root. It is at the first match, or done if there
// are none. If the root or a directory can't be read, or something can't be
// described, the Iterator is done and Err returns the error. Data reads the
// current file, and if that fails the Iterator is done and Err returns the error.
func (mr MatchRoot) Iterator() (i Iterator, done bool) {
	mri := &matchRootIter{
		MatchRoot: mr,
	}
	return mri, mri.Reset()
}

// Factory fulfills liter.Factory, so a MatchRoot can be used by what takes one,
// such as slice.FromIterFactory. The values are the paths.
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
		mri.info, mri.err = CoreFsStat(mri.CoreFS, mri.path)
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
	files, err := ReadDirNames(mri.CoreFS, mri.path)
	if err != nil {
		mri.err, mri.done = err, true
		return
	}
	for _, f := range files {
		mri.files = append(mri.files, path.Join(mri.path, f))
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
		mri.data, mri.err = mri.CoreFS.ReadFile(mri.path)
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
	mri.data = nil
	mri.err = nil
	mri.done = false
	// an empty root is the current directory, and io/fs rejects a trailing slash
	mri.path = path.Clean(mri.Root)
	mri.appendFiles()
	mri.done = mri.done || len(mri.files) == 0
	mri.Next()
	mri.idx = 0

	return mri.done
}
