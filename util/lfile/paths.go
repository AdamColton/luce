package lfile

import (
	"io/ioutil"
	"os"
)

// Paths to be iterated over.
type Paths []string

// Iter fulfills Iterator, iterating over the files in Filenames.
func (fn Paths) Iterator() (Iterator, bool) {
	i := &pathsIterator{
		Paths: fn,
	}
	return i, i.update()
}

// Iterator will iterate over Filenames. If Err is ever not nil, it will stop.
// If Autoload is true, Load will be called when invoking Next. Filename and
// Index indicate the current file.
type pathsIterator struct {
	Paths
	filename string
	done     bool
	Index    int
	data     []byte
	err      error
	info     *os.FileInfo
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

func (i *pathsIterator) Data() []byte {
	if i.data == nil && i.err == nil {
		i.data, i.err = ReadFile(i.filename)
		i.done = i.err != nil
	}
	return i.data
}
func (i *pathsIterator) Err() error {
	return i.err
}

// Next moves to the next file. Returned bool will be true when iteration is
// done.
func (i *pathsIterator) Next() (path string, done bool) {
	i.Index++
	return i.filename, i.update()
}

// Reset the Iter to the start and set the autoload value.
func (i *pathsIterator) Reset() (done bool) {
	i.Index = 0
	i.done = false
	i.data = nil
	return i.update()
}

// ReadFile is a reference to ioutil.ReadFile. It is left exposed for testing.
var ReadFile = ioutil.ReadFile

// Stat is a reference to os.Stat. It is left exposed for testing.
var Stat = os.Stat

func (i *pathsIterator) Stat() (info os.FileInfo) {
	if i.info == nil {
		info, i.err = Stat(i.filename)
		i.info = &info
	} else {
		info = *(i.info)
	}
	return
}

func (i *pathsIterator) update() bool {
	i.done = i.done || i.Index >= len(i.Paths) || i.err != nil
	i.info = nil
	if i.done {
		i.filename = ""
	} else {
		i.filename = i.Paths[i.Index]
		i.data = nil
	}

	return i.done
}
