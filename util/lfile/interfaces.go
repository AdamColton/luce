package lfile

import (
	"io"
	"io/fs"
	"os"
)

type DirReader interface {
	ReadDir(n int) ([]os.DirEntry, error)
}

type DirNameReader interface {
	Readdirnames(n int) (names []string, err error)
}

// Dir is fulfilled by *os.File.
type Dir interface {
	DirReader
	DirNameReader
	Name() string
}

type FSDirReader interface {
	ReadDir(name string) ([]fs.DirEntry, error)
}

type FSLstater interface {
	Lstat(name string) (os.FileInfo, error)
}

type FSOpener interface {
	Open(name string) (fs.File, error)
}

type FSCreator interface {
	Create(name string) (fs.File, error)
}

type FSRemover interface {
	Remove(name string) error
}

type FSStater interface {
	Stat(name string) (os.FileInfo, error)
}

type FSFileReader interface {
	ReadFile(name string) ([]byte, error)
}

type FileLstater interface {
	Lstat() (os.FileInfo, error)
}

type CoreFS interface {
	FSOpener
	FSFileReader
	FSDirReader
}

type FSReader interface {
	CoreFS
	FSStater
	FSOpener
	FSFileReader
}

// File provides an interface fulfilled by *os.File. This allows for testing
// without relying on the actual file system.
type File interface {
	fs.File
	Dir
	io.Writer
}

// Repository is an interface to the file system.
type Repository interface {
	FSOpener
	FSCreator
	FSRemover
	FSStater
	FSLstater
	FSFileReader
	FSDirReader
}
