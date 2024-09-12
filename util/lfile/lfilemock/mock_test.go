package lfilemock_test

import (
	"bytes"
	"io"
	"io/fs"
	"io/ioutil"
	"syscall"
	"testing"
	"time"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/lfile/lfilemock"
	"github.com/adamcolton/luce/util/navigator"
	"github.com/stretchr/testify/assert"
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

	f, err := r.Open("file1.txt")
	assert.NoError(t, err)
	b, err := io.ReadAll(f)
	assert.NoError(t, err)
	assert.Equal(t, file1, string(b))

	f, err = r.Open("dir2")
	assert.NoError(t, err)
	des, err := f.ReadDir(0)
	assert.NoError(t, err)
	slice.Less[fs.DirEntry](func(i, j fs.DirEntry) bool {
		return i.Name() < j.Name()
	}).Sort(des)

	assert.Equal(t, "dir3", des[0].Name())
	assert.True(t, des[0].IsDir())
	assert.Equal(t, fs.FileMode(0), des[0].Type())
	fi, err := des[0].Info()
	assert.NoError(t, err)
	assert.Equal(t, "dir3", fi.Name())
	assert.True(t, fi.IsDir())
	assert.Equal(t, int64(0), fi.Size())
	assert.Equal(t, fs.FileMode(0), fi.Mode())
	assert.Equal(t, time.Time{}, fi.ModTime())
	assert.Nil(t, fi.Sys())

	f, err = r.Open("dir2/dir3")
	assert.NoError(t, err)
	fi2, err := f.Stat()
	assert.NoError(t, err)
	assert.Equal(t, fi, fi2)
	assert.NoError(t, f.Close())

	assert.Equal(t, "file4.txt", des[1].Name())
	assert.False(t, des[1].IsDir())

	f, err = r.Open("dir1/file2.bin")
	assert.NoError(t, err)
	b, err = io.ReadAll(f)
	assert.NoError(t, err)
	assert.Equal(t, file2, b)

	f, err = r.Create("dir1/dir4/file3.txt")
	assert.NoError(t, err)
	assert.Equal(t, "file3.txt", f.Name())

	err = r.Remove("dir1/file2.bin")
	assert.NoError(t, err)
	f, _ = r.Open("dir1/file2.bin")
	assert.Nil(t, f)

	lf := lfilemock.New("test.txt", "this is a test")
	got, err := ioutil.ReadAll(lf)
	assert.NoError(t, err)
	assert.Equal(t, []byte("this is a test"), got)
	te := lerr.Str("test_error")
	lf.Err = te
	_, err = lf.Read(nil)
	assert.Equal(t, te, err)

	data := []byte{3, 1, 4, 1, 5, 9, 2, 6, 5, 3}
	lf = lfilemock.New("test2.txt", data)
	got, err = ioutil.ReadAll(lf)
	assert.NoError(t, err)
	assert.Equal(t, data, got)
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

func TestByteFile(t *testing.T) {
	bf := &lfilemock.ByteFile{
		Name: "Test",
		Data: bytes.NewBuffer([]byte{1, 2, 3, 4, 5}),
	}
	tree, ok := bf.Next("foo", true, navigator.Void)
	assert.Nil(t, tree)
	assert.False(t, ok)
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

	f, err := r.Open("file1.txt")
	assert.NoError(t, err)

	expectErr := &fs.PathError{Op: "readdirent", Path: "file1.txt", Err: syscall.ENOTDIR}
	_, err = f.Readdirnames(-1)
	assert.Equal(t, expectErr, err)

	f, err = r.Open("/")
	assert.NoError(t, err)
	names, err := f.Readdirnames(-1)
	assert.NoError(t, err)
	slice.LT[string]().Sort(names)
	expected := []string{"dir1", "file1.txt", "file2.bin", "file3.txt", "file4.txt"}
	assert.Equal(t, expected, names)
	_, err = f.Readdirnames(-1)
	assert.Equal(t, io.EOF, err)

	f, err = r.Open("/")
	assert.NoError(t, err)
	names, err = f.Readdirnames(2)
	assert.NoError(t, err)
	assert.Len(t, names, 2)
	names, err = f.Readdirnames(-1)
	assert.NoError(t, err)
	assert.Len(t, names, 3)
	_, err = f.Readdirnames(-1)
	assert.Equal(t, io.EOF, err)
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

	f, err := r.Open("/dir1/file2.bin")
	assert.NoError(t, err)
	f.Write([]byte("test"))
	f.Close()

	f2, err := r.Open("/dir1/file2.bin")
	assert.NoError(t, err)
	expected := append(file2, []byte("test")...)
	got, err := io.ReadAll(f2)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)

}
