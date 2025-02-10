package lfilemock_test

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"syscall"
	"testing"
	"time"

	"github.com/adamcolton/luce/ds/lbuf"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/lfile"
	"github.com/adamcolton/luce/util/lfile/lfilemock"
	"github.com/adamcolton/luce/util/navigator"
	"github.com/stretchr/testify/assert"
)

// The mock fulfills the interfaces of lfile, and the files it opens fulfill
// lfile.File.
var (
	_ lfile.Repository = (*lfilemock.Repository)(nil)
	_ lfile.File       = (*lfilemock.File)(nil)
	_ lfile.CoreFS     = (*lfilemock.Repository)(nil)
	_ fs.FS            = (*lfilemock.Repository)(nil)
)

func TestRepository(t *testing.T) {
	file1 := "this is test file 1"
	file2 := []byte{3, 1, 4, 1, 5, 9, 2, 6, 5}
	r := lfilemock.Parse(map[string]any{
		"file1.txt": file1,
		"dir1": map[string]any{
			"file2.bin": file2,
		},
		"dir2": map[string]any{
			"dir3":      map[string]any{},
			"file4.txt": "this is file 4",
		},
	}).Repository()

	f := lerr.Must(r.Open("file1.txt")).(lfile.File)
	b, err := io.ReadAll(f)
	assert.NoError(t, err)
	assert.Equal(t, file1, string(b))

	f = lerr.Must(r.Open("dir2")).(lfile.File)
	des, err := f.ReadDir(0)
	assert.NoError(t, err)
	// the entries are sorted
	assert.Equal(t, "dir3", des[0].Name())
	assert.True(t, des[0].IsDir())
	assert.Equal(t, fs.ModeDir, des[0].Type())
	fi, err := des[0].Info()
	assert.NoError(t, err)
	assert.Equal(t, "dir3", fi.Name())
	assert.True(t, fi.IsDir())
	assert.Equal(t, int64(0), fi.Size())
	assert.Equal(t, fs.ModeDir|0o755, fi.Mode())
	assert.Equal(t, time.Time{}, fi.ModTime())
	assert.Nil(t, fi.Sys())

	f = lerr.Must(r.Open("dir2/dir3")).(lfile.File)
	fi2, err := f.Stat()
	assert.NoError(t, err)
	assert.Equal(t, fi, fi2)
	assert.NoError(t, f.Close())

	assert.Equal(t, "file4.txt", des[1].Name())
	assert.False(t, des[1].IsDir())
	assert.Equal(t, fs.FileMode(0), des[1].Type())

	f = lerr.Must(r.Open("dir1/file2.bin")).(lfile.File)
	b, err = io.ReadAll(f)
	assert.NoError(t, err)
	assert.Equal(t, file2, b)

	f = lerr.Must(r.Create("dir2/dir3/file3.txt")).(lfile.File)
	assert.Equal(t, "dir2/dir3/file3.txt", f.Name(), "a file is named as it was opened")
	_, err = f.Write([]byte("abc"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("abc"), lerr.Must(r.ReadFile("dir2/dir3/file3.txt")))

	err = r.Remove("dir1/file2.bin")
	assert.NoError(t, err)
	nilFile, err := r.Open("dir1/file2.bin")
	assert.Nil(t, nilFile)
	assert.ErrorIs(t, err, fs.ErrNotExist)

	lf := lfilemock.New("test.txt", "this is a test")
	got, err := io.ReadAll(lf)
	assert.NoError(t, err)
	assert.Equal(t, []byte("this is a test"), got)
	te := lerr.Str("test_error")
	lf.Err = te
	_, err = lf.Read(nil)
	assert.Equal(t, te, err)

	data := []byte{3, 1, 4, 1, 5, 9, 2, 6, 5, 3}
	lf = lfilemock.New("test2.txt", data)
	got, err = io.ReadAll(lf)
	assert.NoError(t, err)
	assert.Equal(t, data, got)
}

func TestPaths(t *testing.T) {
	r := lfilemock.Parse(map[string]any{
		"a": map[string]any{"b.txt": "b"},
	}).Repository()

	// a leading or trailing slash, empty parts and "." parts are ignored
	for _, name := range []string{"a/b.txt", "/a/b.txt", "a//b.txt", "a/./b.txt", "a/b.txt/"} {
		f, err := r.Open(name)
		assert.NoError(t, err, name)
		assert.Equal(t, name, f.(lfile.File).Name(), "the name is not changed")
		fi, err := f.Stat()
		assert.NoError(t, err)
		assert.Equal(t, "b.txt", fi.Name())
	}

	for _, name := range []string{"", ".", "/", "//", "./"} {
		des, err := r.ReadDir(name)
		assert.NoError(t, err, name)
		if assert.Len(t, des, 1, name) {
			assert.Equal(t, "a", des[0].Name(), name)
		}
	}
}

func TestNewPanic(t *testing.T) {
	defer func() {
		assert.Equal(t, lfilemock.ErrNewType, recover())
	}()
	lfilemock.New("thisShouldPanic", 123)
}

func TestParseDirPanic(t *testing.T) {
	defer func() {
		assert.Equal(t, lfilemock.ErrParseType, recover())
	}()
	lfilemock.Parse(map[string]any{
		"file1.txt": 123,
	})
}

func TestNamePanic(t *testing.T) {
	bad := []string{"", ".", "..", "a/b"}
	for _, name := range bad {
		for n, fn := range map[string]func(){
			"file":     func() { lfilemock.Parse(map[string]any{name: "x"}) },
			"bytes":    func() { lfilemock.Parse(map[string]any{name: []byte("x")}) },
			"dir":      func() { lfilemock.Parse(map[string]any{name: map[string]any{}}) },
			"files":    func() { lfilemock.Parse(map[string]any{"d": []string{name}}) },
			"strings":  func() { lfilemock.Parse(map[string]any{name: []string{"x"}}) },
			"AddFile":  func() { lfilemock.Parse(nil).AddFile(name, nil) },
			"AddDir":   func() { lfilemock.Parse(nil).AddDir(name) },
			"AddFiles": func() { lfilemock.Parse(nil).AddDir("d").Children["d"].(*lfilemock.Directory).AddFile(name, nil) },
		} {
			if name == "." && n == "strings" {
				// "." is how files are added to the directory itself
				assert.NotPanics(t, fn)
				continue
			}
			assert.PanicsWithValue(t, lfilemock.ErrParseName, fn, "%s %q", n, name)
		}
	}
}

func TestByteFile(t *testing.T) {
	bf := &lfilemock.ByteFile{
		Name: "Test",
		Data: lbuf.New([]byte{1, 2, 3, 4, 5}),
	}
	tree, ok := bf.Next("foo", true, navigator.Void)
	assert.Nil(t, tree)
	assert.False(t, ok)
	assert.NoError(t, bf.Error())
	bf.Err = lerr.Str("boom")
	assert.Equal(t, lerr.Str("boom"), bf.Error())
}

func TestDirectory(t *testing.T) {
	d := lfilemock.Parse(map[string]any{
		"a": map[string]any{"b.txt": "b"},
	})
	n, ok := d.Next("a", false, navigator.Void)
	assert.True(t, ok)
	assert.Equal(t, "a", n.(*lfilemock.Directory).Name)
	_, ok = d.Next("nope", false, navigator.Void)
	assert.False(t, ok)

	// create returns the directory that was added
	n, ok = d.Next("new", true, navigator.Void)
	assert.True(t, ok)
	assert.Equal(t, "new", n.(*lfilemock.Directory).Name)
	assert.Same(t, n, d.Children["new"])

	assert.NoError(t, d.Error())
	d.Err = lerr.Str("boom")
	assert.Equal(t, lerr.Str("boom"), d.Error())

	// AddDir returns the directory it was called on
	assert.Same(t, d, d.AddDir("x"))
	assert.Same(t, d, d.AddDir("y").AddDir("z"))
}

func TestGet(t *testing.T) {
	d := lfilemock.Parse(map[string]any{
		"foo": map[string]any{
			"bar": map[string]any{"x.txt": "x"},
		},
	})

	n, found := d.Get("foo/bar")
	assert.True(t, found)
	assert.Equal(t, "bar", n.(*lfilemock.Directory).Name, "the node is the one in the tree")
	n, found = d.Get("/foo//bar/x.txt")
	assert.True(t, found)
	assert.Equal(t, "x.txt", n.(*lfilemock.ByteFile).Name)
	n, found = d.Get("")
	assert.True(t, found)
	assert.Same(t, d, n)

	for _, name := range []string{"foo/nope", "foo/bar/x.txt/y", "nope"} {
		n, found = d.Get(name)
		assert.False(t, found, name)
		assert.Nil(t, n)
	}
}

func TestReaddirnames(t *testing.T) {
	file1 := "this is test file 1"
	file2 := []byte{3, 1, 4, 1, 5, 9, 2, 6, 5}
	r := lfilemock.Parse(map[string]any{
		"file1.txt": file1,
		"dir1": map[string]any{
			"file2.bin": file2,
		},
		"file2.bin": file2,
		"file3.txt": "file3.txt",
		"file4.txt": "file4.txt",
	}).Repository()

	f := lerr.Must(r.Open("file1.txt")).(lfile.File)

	expectErr := &fs.PathError{Op: "readdirent", Path: "file1.txt", Err: syscall.ENOTDIR}
	_, err := f.Readdirnames(-1)
	assert.Equal(t, expectErr, err)
	_, err = f.ReadDir(-1)
	assert.Equal(t, expectErr, err)

	f = lerr.Must(r.Open("/")).(lfile.File)
	names, err := f.Readdirnames(-1)
	assert.NoError(t, err)
	expected := []string{"dir1", "file1.txt", "file2.bin", "file3.txt", "file4.txt"}
	assert.Equal(t, expected, names, "the names are in order")
	// as os.File does, everything read is not an error when n <= 0
	names, err = f.Readdirnames(-1)
	assert.NoError(t, err)
	assert.Empty(t, names)

	// the position is shared by ReadDir and Readdirnames, and a page that is
	// asked for after the end is io.EOF
	f = lerr.Must(r.Open("/")).(lfile.File)
	names, err = f.Readdirnames(2)
	assert.NoError(t, err)
	assert.Equal(t, expected[:2], names)
	des, err := f.ReadDir(2)
	assert.NoError(t, err)
	assert.Equal(t, "file2.bin", des[0].Name())
	assert.Equal(t, "file3.txt", des[1].Name())
	names, err = f.Readdirnames(10)
	assert.NoError(t, err)
	assert.Equal(t, expected[4:], names)
	_, err = f.Readdirnames(1)
	assert.Equal(t, io.EOF, err)
	_, err = f.ReadDir(1)
	assert.Equal(t, io.EOF, err)

	f.(*lfilemock.File).Err = lerr.Str("boom")
	_, err = f.Readdirnames(-1)
	assert.Equal(t, lerr.Str("boom"), err)
}

func TestWriteByteFile(t *testing.T) {
	file1 := "this is test file 1"
	file2 := []byte{3, 1, 4, 1, 5, 9, 2, 6, 5}
	r := lfilemock.Parse(map[string]any{
		"file1.txt": file1,
		"dir1": map[string]any{
			"file2.bin": file2,
		},
		"dir2": map[string]any{
			"dir3":      map[string]any{},
			"file4.txt": "this is file 4",
		},
	}).Repository()

	f := lerr.Must(r.Open("/dir1/file2.bin")).(lfile.File)
	f.Write([]byte("test"))
	f.Close()

	f2, err := r.Open("/dir1/file2.bin")
	assert.NoError(t, err)
	expected := append(file2, []byte("test")...)
	got, err := io.ReadAll(f2)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)

	// the size Stat reports follows the writes
	fi, err := r.Stat("/dir1/file2.bin")
	assert.NoError(t, err)
	assert.Equal(t, int64(len(expected)), fi.Size())
}

func TestReadTwice(t *testing.T) {
	repo := lfilemock.Parse(map[string]any{
		".": []string{"a", "aa", "b", "c"},
		"dir": map[string]any{
			".":    []string{"d", "e", "f"},
			"dir1": []string{"g", "h", "hh", "i"},
			"dir2": []string{"j", "k"},
		},
		".hidden1": []string{"x", "y", "z"},
	}).Repository()

	f := lerr.Must(repo.Open("/dir/dir1")).(lfile.File)
	names := lerr.Must(f.Readdirnames(-1))
	sort.Strings(names)
	expected := []string{"g", "h", "hh", "i"}
	assert.Equal(t, expected, names)

	f = lerr.Must(repo.Open("/dir/dir1/g")).(lfile.File)
	got := string(lerr.Must(io.ReadAll(f)))
	assert.Equal(t, "g", got)
	f.Close()

	f = lerr.Must(repo.Open("/dir/dir1/g")).(lfile.File)
	got = string(lerr.Must(io.ReadAll(f)))
	assert.Equal(t, "g", got)
}

func TestCreate(t *testing.T) {
	d := lfilemock.Parse(map[string]any{
		"dir":  map[string]any{},
		"file": "hello",
	})
	r := d.Repository()

	// an existing file is emptied, and keeps what was set on it
	when := time.Date(2024, 10, 13, 0, 0, 0, 0, time.UTC)
	n, _ := d.Get("file")
	n.(*lfilemock.ByteFile).Mod = when
	f, err := r.Create("file")
	assert.NoError(t, err)
	fi, err := f.Stat()
	assert.NoError(t, err)
	assert.Equal(t, int64(0), fi.Size())
	assert.Equal(t, when, fi.ModTime())

	f, err = r.Create("/dir/new.txt")
	assert.NoError(t, err)
	_, err = f.(io.Writer).Write([]byte("abc"))
	assert.NoError(t, err)
	fi, err = r.Stat("dir/new.txt")
	assert.NoError(t, err)
	assert.Equal(t, int64(3), fi.Size())
	assert.Equal(t, fs.FileMode(0o644), fi.Mode())

	// the root Err
	d.Err = lerr.Str("boom")
	_, err = r.Create("dir/other")
	assert.Equal(t, d.Err, err)
}

func TestModeModSys(t *testing.T) {
	when := time.Date(2024, 10, 13, 12, 30, 0, 0, time.UTC)
	sys := &syscall.Stat_t{Size: 42}
	d := lfilemock.Parse(map[string]any{
		"dir": map[string]any{"file.txt": "hello"},
	})
	dir, _ := d.Get("dir")
	dir.(*lfilemock.Directory).Mod = when
	dir.(*lfilemock.Directory).Mode = 0o700
	dir.(*lfilemock.Directory).SysData = sys
	file, _ := d.Get("dir/file.txt")
	file.(*lfilemock.ByteFile).Mod = when
	file.(*lfilemock.ByteFile).Mode = 0o600 | fs.ModeDir // ModeDir is not for a file
	file.(*lfilemock.ByteFile).SysData = sys
	r := d.Repository()

	check := func(fi fs.FileInfo, mode fs.FileMode, isDir bool) {
		t.Helper()
		assert.Equal(t, mode, fi.Mode())
		assert.Equal(t, isDir, fi.IsDir())
		assert.Equal(t, isDir, fi.Mode().IsDir())
		assert.Equal(t, when, fi.ModTime())
		assert.Same(t, sys, fi.Sys())
	}

	// by Stat and Lstat, by a File and by the DirEntry that lists it
	for _, stat := range []func(string) (fs.FileInfo, error){r.Stat, r.Lstat} {
		fi, err := stat("dir")
		assert.NoError(t, err)
		check(fi, fs.ModeDir|0o700, true)
		fi, err = stat("dir/file.txt")
		assert.NoError(t, err)
		check(fi, 0o600, false)
	}
	f, _ := r.Open("dir/file.txt")
	fi, err := f.Stat()
	assert.NoError(t, err)
	check(fi, 0o600, false)

	des, err := r.ReadDir("dir")
	assert.NoError(t, err)
	fi, err = des[0].Info()
	assert.NoError(t, err)
	check(fi, 0o600, false)
	assert.Equal(t, fs.FileMode(0), des[0].Type())

	des, err = r.ReadDir("")
	assert.NoError(t, err)
	fi, err = des[0].Info()
	assert.NoError(t, err)
	check(fi, fs.ModeDir|0o700, true)
	assert.Equal(t, fs.ModeDir, des[0].Type())
}

func TestErrorInjection(t *testing.T) {
	boom := lerr.Str("boom")
	d := lfilemock.Parse(map[string]any{
		"dir":  map[string]any{"x": "x"},
		"file": "hello",
	})
	r := d.Repository()
	dir, _ := d.Get("dir")
	dir.(*lfilemock.Directory).Err = boom
	file, _ := d.Get("file")
	file.(*lfilemock.ByteFile).Err = boom

	// the error is on the file that is opened
	f, err := r.Open("file")
	assert.NoError(t, err)
	_, err = f.Read(make([]byte, 1))
	assert.Equal(t, boom, err)
	_, err = f.(io.Writer).Write([]byte("x"))
	assert.Equal(t, boom, err)
	_, err = f.Stat()
	assert.Equal(t, boom, err)
	assert.Equal(t, boom, f.Close())
	_, err = r.Stat("file")
	assert.Equal(t, boom, err)
	_, err = r.ReadFile("file")
	assert.Equal(t, boom, err)
	_, err = r.ReadDir("dir")
	assert.Equal(t, boom, err)

	// the root's error is on the repository
	d.Err = boom
	_, err = r.Open("file")
	assert.Equal(t, boom, err)
	_, err = r.ReadDir("file")
	assert.Equal(t, boom, err)
	assert.Equal(t, boom, r.Remove("file"))
	_, err = r.Lstat("file")
	assert.Equal(t, boom, err)
}

func TestRoot(t *testing.T) {
	r := lfilemock.Parse(map[string]any{"a": "a"}).Repository()

	// the root can't be removed or created
	for _, name := range []string{"", ".", "/"} {
		var pe *fs.PathError
		assert.True(t, errors.As(r.Remove(name), &pe), name)
		assert.Equal(t, syscall.EINVAL, pe.Err)
		_, err := r.Create(name)
		assert.True(t, errors.As(err, &pe), name)
		assert.Equal(t, syscall.EISDIR, pe.Err)
	}
	des, err := r.ReadDir("/")
	assert.NoError(t, err)
	assert.Len(t, des, 1)
}

func TestDirHandle(t *testing.T) {
	r := lfilemock.Parse(map[string]any{"dir": map[string]any{}}).Repository()
	f := lerr.Must(r.Open("dir")).(lfile.File)

	var pe *fs.PathError
	_, err := f.Write([]byte("x"))
	assert.True(t, errors.As(err, &pe))
	assert.Equal(t, "write", pe.Op)
	assert.Equal(t, syscall.EBADF, pe.Err)
	_, err = f.Read(make([]byte, 1))
	assert.True(t, errors.As(err, &pe))
	assert.Equal(t, "read", pe.Op)
	assert.Equal(t, syscall.EISDIR, pe.Err)
}

func TestDirEntryErr(t *testing.T) {
	boom := lerr.Str("boom")
	de := &lfilemock.DirEntry{
		EntryName: "x",
		Err:       boom,
		FileInfo:  &lfilemock.FileInfo{FileName: "x"},
	}
	fi, err := de.Info()
	assert.Equal(t, boom, err)
	assert.Equal(t, "x", fi.Name())
}

// osFS is the os package over a directory, so the mock can be compared with it.
type osFS struct{ root string }

func (o osFS) path(name string) string { return filepath.Join(o.root, name) }

func (o osFS) Open(name string) (fs.File, error) {
	f, err := os.Open(o.path(name))
	if err != nil {
		return nil, err
	}
	return f, nil
}
func (o osFS) Create(name string) (fs.File, error) {
	f, err := os.Create(o.path(name))
	if err != nil {
		return nil, err
	}
	return f, nil
}
func (o osFS) Remove(name string) error               { return os.Remove(o.path(name)) }
func (o osFS) Stat(name string) (fs.FileInfo, error)  { return os.Stat(o.path(name)) }
func (o osFS) Lstat(name string) (fs.FileInfo, error) { return os.Lstat(o.path(name)) }
func (o osFS) ReadFile(name string) ([]byte, error)   { return os.ReadFile(o.path(name)) }
func (o osFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return os.ReadDir(o.path(name))
}

// TestLikeOS runs the same calls on a real directory and on the mock and
// requires the same kind of error: the same operation and the same errno.
func TestLikeOS(t *testing.T) {
	tree := func() (real, mock lfile.Repository) {
		dir := t.TempDir()
		assert.NoError(t, os.MkdirAll(filepath.Join(dir, "dir", "sub"), 0o755))
		assert.NoError(t, os.MkdirAll(filepath.Join(dir, "empty"), 0o755))
		assert.NoError(t, os.WriteFile(filepath.Join(dir, "file.txt"), []byte("hello"), 0o644))
		assert.NoError(t, os.WriteFile(filepath.Join(dir, "dir", "a.txt"), []byte("a"), 0o644))
		mock = lfilemock.Parse(map[string]any{
			"dir":      map[string]any{"sub": map[string]any{}, "a.txt": "a"},
			"empty":    map[string]any{},
			"file.txt": "hello",
		}).Repository()
		return osFS{dir}, mock
	}

	calls := map[string]func(r lfile.Repository) error{
		"Open missing":         func(r lfile.Repository) error { _, err := r.Open("nope"); return err },
		"Open below file":      func(r lfile.Repository) error { _, err := r.Open("file.txt/x"); return err },
		"Open below missing":   func(r lfile.Repository) error { _, err := r.Open("nope/x"); return err },
		"Stat missing":         func(r lfile.Repository) error { _, err := r.Stat("nope"); return err },
		"Stat below file":      func(r lfile.Repository) error { _, err := r.Stat("file.txt/x"); return err },
		"Lstat missing":        func(r lfile.Repository) error { _, err := r.Lstat("nope"); return err },
		"ReadFile missing":     func(r lfile.Repository) error { _, err := r.ReadFile("nope"); return err },
		"ReadFile dir":         func(r lfile.Repository) error { _, err := r.ReadFile("dir"); return err },
		"ReadDir missing":      func(r lfile.Repository) error { _, err := r.ReadDir("nope"); return err },
		"ReadDir file":         func(r lfile.Repository) error { _, err := r.ReadDir("file.txt"); return err },
		"Remove missing":       func(r lfile.Repository) error { return r.Remove("nope") },
		"Remove below missing": func(r lfile.Repository) error { return r.Remove("nope/x") },
		"Remove below file":    func(r lfile.Repository) error { return r.Remove("file.txt/x") },
		"Remove full dir":      func(r lfile.Repository) error { return r.Remove("dir") },
		"Create missing dir":   func(r lfile.Repository) error { _, err := r.Create("nope/x"); return err },
		"Create below file":    func(r lfile.Repository) error { _, err := r.Create("file.txt/x"); return err },
		"Create over dir":      func(r lfile.Repository) error { _, err := r.Create("dir"); return err },
		"Remove file":          func(r lfile.Repository) error { return r.Remove("file.txt") },
		"Remove empty dir":     func(r lfile.Repository) error { return r.Remove("empty") },
		"Create":               func(r lfile.Repository) error { _, err := r.Create("empty/new.txt"); return err },
		"ReadFile":             func(r lfile.Repository) error { _, err := r.ReadFile("file.txt"); return err },
		"ReadDir":              func(r lfile.Repository) error { _, err := r.ReadDir("dir"); return err },
		"Stat dir":             func(r lfile.Repository) error { _, err := r.Stat("dir"); return err },
		"read a dir handle":    func(r lfile.Repository) error { f, _ := r.Open("dir"); _, err := f.Read(make([]byte, 1)); return err },
		"ReadDir of a file handle": func(r lfile.Repository) error {
			f, _ := r.Open("file.txt")
			_, err := f.(lfile.Dir).ReadDir(-1)
			return err
		},
		"Readdirnames of a file": func(r lfile.Repository) error {
			f, _ := r.Open("file.txt")
			_, err := f.(lfile.Dir).Readdirnames(-1)
			return err
		},
	}

	for name, call := range calls {
		t.Run(name, func(t *testing.T) {
			real, mock := tree()
			want, got := call(real), call(mock)
			if want == nil {
				assert.NoError(t, got)
				return
			}
			var wantPE, gotPE *fs.PathError
			if assert.True(t, errors.As(want, &wantPE), "the os error is a PathError: %v", want) &&
				assert.True(t, errors.As(got, &gotPE), "the mock error is a PathError: %v", got) {
				assert.Equal(t, wantPE.Op, gotPE.Op)
				assert.Equal(t, wantPE.Err, gotPE.Err)
			}
		})
	}
}
