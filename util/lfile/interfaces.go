package lfile

import (
	"io"
	"io/fs"
)

// The interfaces in this file describe file systems in small pieces so that a
// function can ask for only what it needs. Most of them are pieces of the
// standard library's io/fs interfaces, and the file systems in the standard
// library fulfill them: *os.File is a Dir and a File, embed.FS and
// testing/fstest.MapFS are a CoreFS. OSRepository wraps the os package to
// fulfill Repository, and the lfilemock package provides a fake to test with.
//
// Paths passed to a file system are slash separated and unrooted, as io/fs
// requires.

// DirReader reads the entries of an open directory. It is the ReadDir method of
// *os.File, so it is not the same as FSDirReader, which is given a path.
type DirReader interface {
	ReadDir(n int) ([]fs.DirEntry, error)
}

// DirNameReader reads the names of the entries of an open directory. It is the
// Readdirnames method of *os.File, which fs.File does not have.
type DirNameReader interface {
	Readdirnames(n int) (names []string, err error)
}

// Dir is an open directory. It is fulfilled by *os.File.
type Dir interface {
	DirReader
	DirNameReader
	Name() string
}

// FSDirReader reads a directory by path. It is the ReadDir half of
// fs.ReadDirFS (that interface also requires Open). It is left without Open so
// a type can offer a cached listing of a directory even when the directory
// can't be read at that moment.
type FSDirReader interface {
	ReadDir(name string) ([]fs.DirEntry, error)
}

// FSLstater describes a file without following a symbolic link. It has no
// counterpart in io/fs.
type FSLstater interface {
	Lstat(name string) (fs.FileInfo, error)
}

// FSOpener opens a file or directory by path. It is fs.FS.
type FSOpener = fs.FS

// FSCreator creates a file, or truncates it if it exists. It has no
// counterpart in io/fs, which is read only.
type FSCreator interface {
	Create(name string) (fs.File, error)
}

// FSRemover removes a file or an empty directory. It has no counterpart in
// io/fs, which is read only.
type FSRemover interface {
	Remove(name string) error
}

// FSStater describes a file by path. It is the Stat half of fs.StatFS (that
// interface also requires Open).
type FSStater interface {
	Stat(name string) (fs.FileInfo, error)
}

// FSFileReader reads a whole file by path. It is the ReadFile half of
// fs.ReadFileFS (that interface also requires Open).
type FSFileReader interface {
	ReadFile(name string) ([]byte, error)
}

// FileLstater is an optional method of an open file. It describes the file
// without following a symbolic link. WrapLstat uses it when the file system has
// no Lstat.
type FileLstater interface {
	Lstat() (fs.FileInfo, error)
}

// CoreFS is the least a file system needs to be searched and read. It is what
// embed.FS and os.DirFS provide, and it is the same set of methods as
// fs.ReadFileFS and fs.ReadDirFS together.
type CoreFS interface {
	FSOpener
	FSFileReader
	FSDirReader
}

// FSReader is a CoreFS that can also describe a file by path. It is what
// fs.StatFS adds to a CoreFS. WrapCoreFS makes an FSReader out of a CoreFS.
type FSReader interface {
	CoreFS
	FSStater
}

// File is an open file that is also a Dir, so it can be read, listed and
// written. It is fulfilled by *os.File. This allows for testing without relying
// on the actual file system.
type File interface {
	fs.File
	Dir
	io.Writer
}

// Repository is a file system that can be read and changed. It is fulfilled by
// OSRepository and by the mock in lfilemock.
type Repository interface {
	FSOpener
	FSCreator
	FSRemover
	FSStater
	FSLstater
	FSFileReader
	FSDirReader
}
