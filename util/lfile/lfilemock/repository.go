package lfilemock

import (
	"io"
	"io/fs"
	"os"
	"strings"
	"syscall"

	"github.com/adamcolton/luce/util/lfile"
	"github.com/adamcolton/luce/util/navigator"
)

// Repository fulfills lfile.Repository and is intended to mock a File system.
// Paths are slash separated. A leading slash, empty parts and "." parts are
// ignored, so "/dir//file" and "dir/./file" are "dir/file", and "", "." and "/"
// are the root. Errors are *fs.PathError, as the os package returns.
type Repository struct {
	Node
}

func pathErr(op, name string, err error) error {
	return &fs.PathError{Op: op, Path: name, Err: err}
}

// pathParts splits a path into the names it passes through.
func pathParts(name string) []string {
	parts := strings.Split(name, "/")
	out := parts[:0]
	for _, p := range parts {
		if p != "" && p != "." {
			out = append(out, p)
		}
	}
	return out
}

// lookup finds the node at name below root. A part that is not there is
// syscall.ENOENT, and one that is inside a file is syscall.ENOTDIR.
func lookup(root Node, op, name string) (Node, error) {
	cur := root
	for _, part := range pathParts(name) {
		d, ok := cur.(*Directory)
		if !ok {
			return nil, pathErr(op, name, syscall.ENOTDIR)
		}
		next, found := d.Children[part]
		if !found {
			return nil, pathErr(op, name, syscall.ENOENT)
		}
		cur = next
	}
	return cur, nil
}

// open finds the node and creates the File for it, named as the caller gave it.
func (r *Repository) open(op, name string) (*File, error) {
	if err := r.Error(); err != nil {
		return nil, err
	}
	n, err := lookup(r.Node, op, name)
	if err != nil {
		return nil, err
	}
	f := n.File()
	f.FileName = name
	return f, nil
}

// Open fulfills lfile.Repository. It opens a file or directory. The File is a
// *File, which fulfills lfile.File, and its name is the name given here. A path
// that does not exist is an error for which errors.Is(err, fs.ErrNotExist) is
// true.
func (r *Repository) Open(name string) (fs.File, error) {
	f, err := r.open("open", name)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// ReadDir fulfills lfile.Repository. It lists the directory sorted by name. A
// file is an error, as it is for os.ReadDir, which opens the path as a
// directory.
func (r *Repository) ReadDir(name string) ([]fs.DirEntry, error) {
	f, err := r.open("open", name)
	if err != nil {
		return nil, err
	}
	if !f.Dir {
		return nil, pathErr("open", name, syscall.ENOTDIR)
	}
	return f.ReadDir(-1)
}

// Remove fulfills lfile.Repository. It removes a file or an empty directory.
// Removing what is not there, a directory that has something in it or the root
// is an error.
func (r *Repository) Remove(name string) error {
	if err := r.Error(); err != nil {
		return err
	}
	parts := pathParts(name)
	if len(parts) == 0 {
		return pathErr("remove", name, syscall.EINVAL)
	}
	last := len(parts) - 1
	parent, err := lookup(r.Node, "remove", strings.Join(parts[:last], "/"))
	if err != nil {
		return pathErr("remove", name, err.(*fs.PathError).Err)
	}
	d, ok := parent.(*Directory)
	if !ok {
		return pathErr("remove", name, syscall.ENOTDIR)
	}
	child, found := d.Children[parts[last]]
	if !found {
		return pathErr("remove", name, syscall.ENOENT)
	}
	if cd, ok := child.(*Directory); ok && len(cd.Children) > 0 {
		return pathErr("remove", name, syscall.ENOTEMPTY)
	}
	delete(d.Children, parts[last])
	return nil
}

// Create fulfills lfile.Repository. It creates an empty file, or empties the
// file that is there. The directory it goes in must exist. The File is a *File,
// so it can be written to, and what is written is kept in the tree.
func (r *Repository) Create(name string) (fs.File, error) {
	if err := r.Error(); err != nil {
		return nil, err
	}
	parts := pathParts(name)
	if len(parts) == 0 {
		return nil, pathErr("open", name, syscall.EISDIR)
	}
	last := len(parts) - 1
	parent, err := lookup(r.Node, "open", strings.Join(parts[:last], "/"))
	if err != nil {
		return nil, pathErr("open", name, err.(*fs.PathError).Err)
	}
	d, ok := parent.(*Directory)
	if !ok {
		return nil, pathErr("open", name, syscall.ENOTDIR)
	}

	var bf *ByteFile
	switch existing := d.Children[parts[last]].(type) {
	case *Directory:
		return nil, pathErr("open", name, syscall.EISDIR)
	case *ByteFile:
		bf = existing
		bf.Data.Data, bf.Data.Idx = nil, 0
	default:
		bf = d.AddFile(parts[last], nil)
	}
	f := bf.File()
	f.FileName = name
	return f, nil
}

// Stat fulfills lfile.Repository.
func (r *Repository) Stat(name string) (os.FileInfo, error) {
	return r.stat("stat", name)
}

// Lstat fulfills lfile.Repository. There are no links, so it is the same as
// Stat.
func (r *Repository) Lstat(name string) (os.FileInfo, error) {
	return r.stat("lstat", name)
}

func (r *Repository) stat(op, name string) (os.FileInfo, error) {
	f, err := r.open(op, name)
	if err != nil {
		return nil, err
	}
	return f.Stat()
}

// ReadFile fulfills lfile.Repository. Reading a directory is an error.
func (r *Repository) ReadFile(name string) ([]byte, error) {
	f, err := r.open("open", name)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(f)
}

// Node is fulfilled by ByteFile and Directory creating a mock directory tree.
type Node interface {
	// File creates a File to open the node with.
	File() *File
	// DirEntry describes the node as an entry of its parent.
	DirEntry() os.DirEntry
	// Next finds a child, and creates it if create is true.
	Next(key string, create bool, vc navigator.VoidContext) (Node, bool)
	// Error is the error that is returned by the operations of the node.
	Error() error
}

var _ lfile.Repository = (*Repository)(nil)
