package lfilemock

import (
	"bytes"
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
	f, _ := r.navigator(getPath(name)).Seek(false, navigator.Void)
	if f == nil {
		return nil, nil
	}
	return f.File(), nil
}

// Remove fulfills lfile.Repository. It removes a file if it exists.
func (r *Repository) Remove(name string) error {
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

// ByteFile is used to create files in mock directory trees. ByteFiles are
// created when calling ParseDir.
type ByteFile struct {
	Name string
	Data []byte
}

// File fulfills Node. It uses a ByteFile to create an instance of File that
// fulfills lfile.File.
func (f *ByteFile) File() *File {
	return &File{
		FileName: f.Name,
		Buffer:   bytes.NewBuffer(f.Data),
		Dir:      false,
		FileSize: int64(len(f.Data)),
		FileInfo: &FileInfo{
			FileName: f.Name,
			Dir:      false,
		},
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

// Node is fulfilled by ByteFile and Directory creating a mock directory tree.
type Node interface {
	File() *File
	DirEntry() os.DirEntry
	Next(key string, create bool, vc navigator.VoidContext) (Node, bool)
}

// Directory is used to create mock directory trees. Directories are
// created when calling ParseDir.
type Directory struct {
	Name     string
	Children map[string]Node
	Err      error
}

// DirEntry fulfills Node. It creates a DirEntry instance for the Directory.
func (d *Directory) DirEntry() os.DirEntry {
	return &DirEntry{
		EntryName: d.Name,
		Dir:       true,
		FileInfo: &FileInfo{
			FileName: d.Name,
			Dir:      true,
		},
	}
}

// File fulfills Node. It creates an instance of *File that fulfills lfile.File.
func (d *Directory) File() *File {
	des := make([]os.DirEntry, 0, len(d.Children))
	for _, c := range d.Children {
		des = append(des, c.DirEntry())
	}
	return &File{
		FileName:   d.Name,
		Dir:        true,
		DirEntries: des,
		FileInfo: &FileInfo{
			FileName: d.Name,
			Dir:      true,
		},
		FileMode: 0,
		Err:      d.Err,
	}
}

// Next fulfills Node and navigator.Nexter. It is used internally for navigating
// the mock directory.
func (d *Directory) Next(key string, create bool, _ navigator.VoidContext) (Node, bool) {
	n, found := d.Children[key]
	if !found && create {
		n = d.AddDir(key)
		found = true
	}
	return n, found
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

func (d *Directory) Get(path string) (n Node, found bool) {
	vc := navigator.VoidContext{}
	n = d
	for _, key := range strings.Split(path, "/") {
		n, found = n.Next(key, false, vc)
		if !found {
			break
		}
	}
	if found {
		n = pathNode{
			name: path,
			Node: n,
		}
	}
	return
}

// AddDir adds a sub directory. This can be used when modifying the mock
// directory tree.
func (d *Directory) AddDir(name string) *Directory {
	c := &Directory{
		Name:     name,
		Children: make(map[string]Node),
	}
	d.Children[name] = c
	return d
}

// AddFile adds a ByteFile. This can be used when modifying the mock directory
// tree.
func (d *Directory) AddFile(name string, contents []byte) *ByteFile {
	f := &ByteFile{
		Name: name,
		Data: contents,
	}
	d.Children[name] = f
	return f
}

// Parse allows for mock directory trees to be setup easily. To create a
// File the value can be either a string or a []byte. To create a sub directory
// the value should be another map[string]any which will be recursivly parsed.
// The returned Repository fulfills lfile.Repository.
func (d *Directory) Repository() lfile.Repository {
	return &Repository{
		Node: d,
	}
}

// ParseDir allows for mock directory trees to be setup easily. To create a File
// the value can be either a string or a []byte. To create a sub directory the
// value should be another map[string]any which will be recursivly parsed.
// Adding a []string will create a directory where each file's name and contents
// are the same.
func Parse(root map[string]any) *Directory {
	out := &Directory{
		Children: make(map[string]Node),
	}
	for name, f := range root {
		switch x := f.(type) {
		case map[string]any:
			s := Parse(x)
			s.Name = name
			out.Children[name] = s
		case string:
			out.AddFile(name, []byte(x))
		case []string:
			var s *Directory
			if name == "." {
				s = out
			} else {
				s = &Directory{
					Name:     name,
					Children: make(map[string]Node, len(x)),
				}
				out.Children[name] = s
			}
			for _, file := range x {
				s.AddFile(file, []byte(file))
			}
		case []byte:
			out.AddFile(name, x)
		default:
			panic(ErrParseType)
		}
	}
	return out
}
