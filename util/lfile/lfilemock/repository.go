package lfilemock

import (
	"io"
	"io/fs"
	"os"
	"strings"
	"syscall"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/lfile"
	"github.com/adamcolton/luce/util/navigator"
)

// Repository fulfills lfile.Repository and is intended to mock a File system.
type Repository struct {
	Node
}

func (r *Repository) navigator(path []string) *navigator.Navigator[string, Node, navigator.VoidContext] {
	return &navigator.Navigator[string, Node, navigator.VoidContext]{
		Cur:  r.Node,
		Keys: path,
	}
}

var notEmptyStr = filter.NEQ("")

func getPath(name string) slice.Slice[string] {
	return notEmptyStr.Slice(strings.Split(name, "/"))
}

// Open fulfills lfile.Repository. It opens a file or directory if it exists.
func (r *Repository) Open(name string) (lfile.File, error) {
	if err := r.Error(); err != nil {
		return nil, err
	}
	f, _ := r.navigator(getPath(name)).Seek(false, navigator.Void)
	if f == nil {
		return nil, nil
	}
	return f.File(), nil
}

// Remove fulfills lfile.Repository. It removes a file if it exists.
func (r *Repository) Remove(name string) error {
	if err := r.Error(); err != nil {
		return err
	}
	file, path := getPath(name).Pop()
	n := r.navigator(path).Trace(true)
	f, _ := n.Seek(false, navigator.Void)
	if f != nil {
		d, ok := n.Pop()
		if ok {
			d := d.(*Directory)
			delete(d.Children, file)
		}
	}
	return nil
}

// Create fulfills lfile.Repository. It creates a file.
func (r *Repository) Create(name string) (lfile.File, error) {
	if err := r.Error(); err != nil {
		return nil, err
	}
	path := getPath(name)
	name, path = path.Pop()
	t, _ := r.
		navigator(path).
		Seek(true, navigator.Void)
	f := t.(*Directory).
		AddFile(name, nil).
		File()
	return f, nil
}

func (r *Repository) Stat(name string) (os.FileInfo, error) {
	f, err := r.Open(name)
	if err != nil {
		return nil, err
	}
	if f == nil {
		return nil, &fs.PathError{Op: "stat", Path: name, Err: syscall.ENOENT}
	}
	return f.Stat()
}

func (r *Repository) Lstat(name string) (os.FileInfo, error) {
	return r.Stat(name)
}

func (r *Repository) ReadFile(name string) ([]byte, error) {
	f, err := r.Open(name)
	if err != nil {
		return nil, err
	}
	if f == nil {
		return nil, &fs.PathError{Op: "readfile", Path: name, Err: syscall.ENOENT}
	}
	return io.ReadAll(f)
}

// Node is fulfilled by ByteFile and Directory creating a mock directory tree.
type Node interface {
	File() *File
	DirEntry() os.DirEntry
	Next(key string, create bool, vc navigator.VoidContext) (Node, bool)
	Error() error
}

// TODO: This seems like a terrible hack indicative of a larger issue with
// not tracking the path correctly
type pathNode struct {
	name string
	Node
}

func (pn pathNode) File() *File {
	f := pn.Node.File()
	f.FileName = pn.name
	return f
}
