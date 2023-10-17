package lfile

import (
	"os"

	"github.com/adamcolton/luce/util/liter"
)

// IteratorSource can generate an Iterator.
type IteratorSource interface {
	// Iterator returns an Iterator that is at its first value, and true if it
	// has none.
	Iterator() (i Iterator, done bool)
}

// Iterator over a set of files and directories. It is a liter.Iter of the
// paths, and adds what is known about the current path. Next moves to the next
// path and returns it, and true when there are no more. Once an error has been
// found the Iterator is done and Err returns it.
type Iterator interface {
	liter.Iter[string]

	// Path to the current file or directory including the name
	Path() string
	// Data is the contents of the current file. It is read the first time it is
	// asked for, and if it can't be read Err is set and the Iterator is done.
	Data() []byte
	// Err is the error that ended the iteration, if any.
	Err() error
	// Stat describes the current file or directory. It means nothing once the
	// Iterator is done.
	Stat() os.FileInfo
	// Reset moves the Iterator back to its first value and returns true if it has
	// none.
	Reset() (done bool)
}
