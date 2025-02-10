package lfile

import (
	"io/fs"
	"os"
	"path/filepath"
)

// OSRepository fulfills Repository by using functions from the "os" package.
// Paths are slash separated, as they are everywhere in lfile, and are converted
// to the separator of the operating system before they reach os. OSRepository
// is the only place that conversion happens.
type OSRepository struct{}

// Open is a wrapper for os.Open. The file is an *os.File.
func (OSRepository) Open(name string) (fs.File, error) {
	f, err := os.Open(filepath.FromSlash(name))
	if err != nil {
		// returning f would be a non-nil fs.File holding a nil *os.File
		return nil, err
	}
	return f, nil
}

// Create is a wrapper for os.Create. The file is an *os.File, so it can be
// written to.
func (OSRepository) Create(name string) (fs.File, error) {
	f, err := os.Create(filepath.FromSlash(name))
	if err != nil {
		return nil, err
	}
	return f, nil
}

// Remove is a wrapper for os.Remove
func (OSRepository) Remove(name string) error {
	return os.Remove(filepath.FromSlash(name))
}

// Stat is a wrapper for os.Stat. It follows symbolic links.
func (OSRepository) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(filepath.FromSlash(name))
}

// Lstat is a wrapper for os.Lstat. It describes a symbolic link, not its
// target.
func (OSRepository) Lstat(name string) (fs.FileInfo, error) {
	return os.Lstat(filepath.FromSlash(name))
}

// ReadFile is a wrapper for os.ReadFile
func (OSRepository) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(filepath.FromSlash(name))
}

// ReadDir is a wrapper for os.ReadDir. The entries are sorted by name.
func (OSRepository) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(filepath.FromSlash(name))
}
