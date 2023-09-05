package fileservice

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/adamcolton/luce/tools/server/service"
	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/lfile"
)

// Extension represents a file extension.
type Extension string

// Service maps a URL to a directory and provides handlers for the files.
type Service struct {
	BaseURL   string
	RootDir   string
	HandleDir Handler
	Handlers  map[Extension]Handler
	parent    *Services
}

var repo lfile.Repository = lfile.OSRepository{}

func (s *Service) Handler(req *service.Request) *service.Response {
	path := strings.Replace(req.Path, s.BaseURL, s.RootDir, 1)
	if strings.Contains(path, "//") {
		fmt.Println("FIX ME")
		path = strings.ReplaceAll(path, "//", "/")
	}
	resp := req.Response(nil)

	f, err := repo.Open(path)
	if resp.ErrCheck(err) {
		return resp
	}

	info, err := f.Stat()
	if resp.ErrCheck(err) {
		return resp
	}

	h := s.getHandler(info.IsDir(), path)
	return h(path, f, req)
}

func (s *Service) getHandler(isDir bool, path string) Handler {
	if isDir {
		return s.getDirHandler()
	}
	return s.getFileHandler(path)
}

var nilHandlerFilter = filter.New(func(h Handler) bool {
	return h != nil
})

func (s *Service) getDirHandler() Handler {
	var ph Handler
	if s.parent != nil {
		ph = s.parent.DirHandler
	}
	h, _ := nilHandlerFilter.First(s.HandleDir, ph, HandleDir)
	return h
}

func (s *Service) getFileHandler(path string) Handler {
	ext := Extension(filepath.Ext(path))
	if s.Handlers != nil {
		if h, found := s.Handlers[ext]; found {
			return h
		}
	}
	h, found := s.parent.Handlers[ext]
	if found {
		return h
	}
	return HandleBinary
}

func (s *Service) Register(conn *service.Client) {
	conn.Add(
		s.Handler,
		service.NewRoute(s.BaseURL).Get().WithPrefix(),
	)
}
