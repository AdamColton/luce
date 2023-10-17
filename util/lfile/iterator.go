package lfile

import (
	"os"

	"github.com/adamcolton/luce/util/liter"
)

// IteratorSource can generate an Iterator.
type IteratorSource interface {
	Iterator() (i Iterator, done bool)
}

// Iterator over a set of files and directories.
type Iterator interface {
	liter.Iter[string]

	// Path to the current file or directory including the name
	Path() string
	Data() []byte
	Err() error
	Stat() os.FileInfo
	Reset() (done bool)
}
