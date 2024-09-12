package lfilemock

import (
	"bytes"
	"os"

	"github.com/adamcolton/luce/util/navigator"
)

// ByteFile is used to create files in mock directory trees. ByteFiles are
// created when calling ParseDir.
type ByteFile struct {
	Name string
	Data *bytes.Buffer
	Err  error
}

// File fulfills Node. It uses a ByteFile to create an instance of File that
// fulfills lfile.File.
func (f *ByteFile) File() *File {
	f.Data = bytes.NewBuffer(f.Data.Bytes())
	return &File{
		FileName: f.Name,
		Buffer:   f.Data,
		Dir:      false,
		FileSize: int64(f.Data.Len()),
		FileInfo: &FileInfo{
			FileName: f.Name,
			Dir:      false,
		},
		Err: f.Err,
	}
}

// DirEntry fulfills Node. It creates a DirEntry instance for the ByteFile.
func (f *ByteFile) DirEntry() os.DirEntry {
	return &DirEntry{
		EntryName: f.Name,
		Dir:       false,
		FileInfo: &FileInfo{
			FileName: f.Name,
			Dir:      false,
		},
	}
}

// Next fulfills Node and navigator.Nexter. It is used internally for navigating
// the mock directory.
func (f *ByteFile) Next(key string, create bool, _ navigator.VoidContext) (Node, bool) {
	return nil, false
}

func (f *ByteFile) Error() error {
	return f.Err
}
