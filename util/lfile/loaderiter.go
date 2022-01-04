package lfile

import "io/ioutil"

// Iterator is used to iterate over a set of files.
type Iterator interface {
	Iter(autoload bool) (i *Iter, done bool)
}

// Filenames to be iterated over.
type Filenames []string

// Iter fulfills Iterator, iterating over the files in Filenames.
func (fn Filenames) Iter(autoload bool) (i *Iter, done bool) {
	i = &Iter{
		Filenames: fn,
		Autoload:  autoload,
	}
	return i, i.update()
}

// Iter will iterate over Filenames. If Err is ever not nil, it will stop.
// If Autoload is true, Load will be called when invoking Next. Filename and
// Index indicate the current file.
type Iter struct {
	Filenames
	Filename       string
	Done, Autoload bool
	Index          int
	Data           []byte
	Err            error
}

// Next moves to the next file. Returned bool will be true when iteration is
// done.
func (i *Iter) Next() (done bool) {
	i.Index++
	return i.update()
}

// ReadFile is a reference to ioutil.ReadFile. It is left exposed for testing.
var ReadFile = ioutil.ReadFile

// Load the current file to Data. Any errors will be stored in Err.
func (i *Iter) Load() {
	i.Data, i.Err = ReadFile(i.Filename)
	i.Done = i.Err != nil
}

func (i *Iter) update() bool {
	i.Done = i.Done || i.Index >= len(i.Filenames) || i.Err != nil
	if i.Done {
		i.Filename = ""
	} else {
		i.Filename = i.Filenames[i.Index]
		if i.Autoload {
			i.Load()
		}
	}

	return i.Done
}
