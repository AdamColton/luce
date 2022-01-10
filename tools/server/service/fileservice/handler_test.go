package fileservice

import (
	"encoding/json"
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/tools/server/service"
	"github.com/adamcolton/luce/util/lfile"
	"github.com/adamcolton/luce/util/lfile/lfilemock"
	"github.com/stretchr/testify/assert"
)

type statusErr struct {
	error
	status int
}

func (se statusErr) Status() int {
	return se.status
}

func TestByteHandler(t *testing.T) {
	f := lfilemock.New("test", "this is a test")
	req := &service.Request{
		ID: 123,
	}
	resp := HandleBinary("foo", f, req)

	assert.Equal(t, req.ID, resp.ID)
	assert.Equal(t, "this is a test", string(resp.Body))

	f.Err = statusErr{
		error:  lerr.Str("Test Error"),
		status: 543,
	}
	f.Buffer.Reset()

	resp = HandleBinary("foo", f, req)
	assert.Equal(t, "Test Error", string(resp.Body))
	assert.Equal(t, 543, resp.Status)
}

func TestJsonDir(t *testing.T) {
	d := lfilemock.ParseDir(map[string]any{
		"foo": "",
		"baz": map[string]any{},
	})
	d.Name = "testdir"
	f := d.File()
	req := &service.Request{
		ID: 123,
	}
	resp := HandleDir("bar", f, req)
	dc := &lfile.DirContents{}
	json.Unmarshal(resp.Body, dc)
	expected := &lfile.DirContents{
		Name:    "testdir",
		Path:    "bar",
		SubDirs: []string{"baz"},
		Files:   []string{"foo"},
	}
	assert.Equal(t, expected, dc)
}
