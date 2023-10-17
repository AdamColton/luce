package lfile

import (
	"os"

	"github.com/adamcolton/luce/util/iter"
)

// IteratorSource can generate an Iterator.
type IteratorSource interface {
	Iterator() (i Iterator, done bool)
}

// Iterator over a set of files and directories.
type Iterator interface {
	iter.Iter[string]
	Path() string
	Data() []byte
	Err() error
	Stat() os.FileInfo
	Reset() (done bool)
}
