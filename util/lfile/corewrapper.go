package lfile

import (
	"io/fs"
	"os"
)

type coreWrapper struct {
	CoreFS
}

func (ew coreWrapper) Stat(name string) (os.FileInfo, error) {
	f, err := ew.CoreFS.Open(name)
	if err != nil {
		return nil, err
	}
	return f.Stat()
}

func (ew coreWrapper) Open(name string) (fs.File, error) {
	f, err := ew.CoreFS.Open(name)
	if err != nil {
		return nil, err
	}
	out := &fsFileWrapper{
		File: f,
		name: name,
		ew:   ew,
	}
	return out, err
}

func WrapCoreFS(fs CoreFS) FSReader {
	return coreWrapper{fs}
}

type fsFileWrapper struct {
	fs.File
	name string
	ew   coreWrapper
}

func (f *fsFileWrapper) Name() string {
	return f.name
}

func (f *fsFileWrapper) ReadDir(n int) ([]os.DirEntry, error) {
	return f.ew.ReadDir(f.name)
}

func (f *fsFileWrapper) Readdirnames(n int) (names []string, err error) {
	des, err := f.ReadDir(n)
	if err != nil {
		return nil, err
	}
	out := make([]string, len(des))
	for i, de := range des {
		out[i] = de.Name()
	}
	return out, nil
}
