package lfile

import "os"

// IteratorSource can generate an Iterator.
type IteratorSource interface {
	Iterator() (i Iterator, done bool)
}

// Iterator over a set of files and directories.
type Iterator interface {
	Path() string
	Done() bool
	Data() []byte
	Err() error
	Next() (done bool)
	Stat() os.FileInfo
	Reset() (done bool)
}
