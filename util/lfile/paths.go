package lfile

import "io/ioutil"

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
}

func (i *pathsIterator) Path() string {
	return i.filename
}
func (i *pathsIterator) Done() bool {
	return i.done
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
func (i *pathsIterator) Next() (done bool) {
	i.Index++
	return i.update()
}

// ReadFile is a reference to ioutil.ReadFile. It is left exposed for testing.
var ReadFile = ioutil.ReadFile

func (i *pathsIterator) update() bool {
	i.done = i.done || i.Index >= len(i.Paths) || i.err != nil
	if i.done {
		i.filename = ""
	} else {
		i.filename = i.Paths[i.Index]
		i.data = nil
	}

	return i.done
}
