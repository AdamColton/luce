package lfile

import (
	"io/fs"
	"os"
)

// Paths is a list of files to be iterated over. It fulfills IteratorSource, and
// the files are read from the operating system. Use FS to read them from another
// file system.
type Paths []string

// Iterator fulfills IteratorSource, iterating over the files in Paths.
func (fn Paths) Iterator() (Iterator, bool) {
	return PathsFS{Paths: fn}.Iterator()
}

// FS returns an IteratorSource over Paths that reads and describes the files
// with cfs, for example an embed.FS or a mock from lfilemock.
func (fn Paths) FS(cfs CoreFS) PathsFS {
	return PathsFS{
		Paths:  fn,
		CoreFS: cfs,
	}
}

// PathsFS is a list of files to be iterated over, read from a CoreFS. It
// fulfills IteratorSource. If CoreFS is nil, OSRepository is used.
type PathsFS struct {
	Paths  Paths
	CoreFS CoreFS
}

// Iterator fulfills IteratorSource, iterating over the files in Paths.
func (p PathsFS) Iterator() (Iterator, bool) {
	i := &pathsIterator{
		paths: p.Paths,
		cfs:   p.CoreFS,
	}
	if i.cfs == nil {
		i.cfs = OSRepository{}
	}
	return i, i.update()
}

// pathsIterator iterates over paths. If err is ever not nil, it will stop. Index
// is the position of the current path. The data and info of the current path
// are read when they are first asked for.
type pathsIterator struct {
	paths    Paths
	cfs      CoreFS
	filename string
	done     bool
	Index    int
	data     []byte
	err      error
	info     fs.FileInfo
	statted  bool
}

func (i *pathsIterator) Path() string {
	return i.filename
}
func (i *pathsIterator) Done() bool {
	return i.done
}
func (i *pathsIterator) Cur() (path string, done bool) {
	return i.filename, i.done
}

func (i *pathsIterator) Idx() int {
	return i.Index
}

// Data reads the current file. If it can't be read the error is kept and the
// iteration is done.
func (i *pathsIterator) Data() []byte {
	if i.data == nil && i.err == nil {
		i.data, i.err = i.cfs.ReadFile(i.filename)
		i.done = i.err != nil
	}
	return i.data
}
func (i *pathsIterator) Err() error {
	return i.err
}

// Next moves to the next file. It returns the path of that file and true when
// iteration is done.
func (i *pathsIterator) Next() (path string, done bool) {
	i.Index++
	done = i.update()
	return i.filename, done
}

// Reset the iterator to the start, forgetting any error. It returns true when
// iteration is done.
func (i *pathsIterator) Reset() (done bool) {
	i.Index = 0
	i.done = false
	i.err = nil
	i.data = nil
	return i.update()
}

// Stat describes the current file. If it can't be described the error is kept
// and the iteration ends at the next file.
func (i *pathsIterator) Stat() os.FileInfo {
	if !i.statted {
		i.info, i.statted = nil, true
		var err error
		if i.info, err = fs.Stat(i.cfs, i.filename); err != nil {
			i.err = err
		}
	}
	return i.info
}

func (i *pathsIterator) update() bool {
	i.done = i.done || i.Index >= len(i.paths) || i.err != nil
	i.info, i.statted = nil, false
	if i.done {
		i.filename = ""
	} else {
		i.filename = i.paths[i.Index]
		i.data = nil
	}

	return i.done
}
