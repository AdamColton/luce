package fileservice

import (
	"testing"

	"github.com/adamcolton/luce/tools/server/service"
	"github.com/adamcolton/luce/util/lfile"
	"github.com/adamcolton/luce/util/lfile/lfilemock"
	"github.com/stretchr/testify/assert"
)

func TestServiceGetHandler(t *testing.T) {
	s := &Service{
		BaseURL: "/",
	}
	conn := &service.Client{
		Mux: service.NewMux(),
	}

	s.Register(conn)

	isHandler := false
	runTest := func(isDir bool, path string) {
		isHandler = false
		h := s.getHandler(isDir, path)
		h("", nil, nil)
		assert.True(t, isHandler)
	}

	restoreHandleDir := HandleDir
	HandleDir = func(path string, f lfile.File, req *service.Request) *service.Response {
		isHandler = true
		return nil
	}
	runTest(true, "/")

	restoreHandleBinary := HandleBinary
	HandleBinary = HandleDir
	s.HandleDir = HandleDir
	HandleDir = restoreHandleDir
	runTest(true, "/")

	s.Handlers = map[Extension]Handler{
		".txt": HandleBinary,
	}
	runTest(false, "foo.txt")

	s.parent = &Services{
		Handlers: map[Extension]Handler{
			".bar": s.Handlers[".txt"],
		},
	}
	delete(s.Handlers, ".txt")
	runTest(false, "foo.bar")

	s.HandleDir = nil
	runTest(false, "foo.baz")
	HandleBinary = restoreHandleBinary

}

func TestServiceHandler(t *testing.T) {
	s := &Service{
		BaseURL: "/",
		parent: &Services{
			Handlers: map[Extension]Handler{},
		},
	}
	conn := &service.Client{
		Mux: service.NewMux(),
	}

	s.Register(conn)

	name := "foo.txt"
	text := "this is a test"
	repo = lfilemock.Parse(map[string]any{
		name: text,
	})

	req := &service.Request{
		ID:   123,
		Path: name,
	}
	resp := s.Handler(req)
	assert.Equal(t, req.ID, resp.ID)
	assert.Equal(t, text, string(resp.Body))
}
