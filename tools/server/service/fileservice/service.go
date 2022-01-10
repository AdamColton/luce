package fileservice

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/adamcolton/luce/tools/server/service"
)

type Service struct {
	BaseURL   string
	RootDir   string
	HandleDir Handler
	Handlers  map[string]Handler
	parent    *Services
}

func (s *Service) Handler(req *service.Request) *service.Response {
	path := strings.Replace(req.Path, s.BaseURL, s.RootDir, 1)
	if strings.Contains(path, "//") {
		fmt.Println("FIX ME")
		path = strings.ReplaceAll(path, "//", "/")
	}

	f, err := os.Open(path)
	if err != nil {
		return req.ResponseErr(err, 404)
	}

	info, err := f.Stat()
	if err != nil {
		return req.ResponseErr(err, 500)
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

func (s *Service) getDirHandler() Handler {
	if s.HandleDir != nil {
		return s.HandleDir
	}
	return HandleDir
}

func (s *Service) getFileHandler(path string) Handler {
	ext := filepath.Ext(path)
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

type Services struct {
	Handlers map[string]Handler
	byURL    map[string]*Service
}

func New() *Services {
	return &Services{
		Handlers: make(map[string]Handler),
		byURL:    make(map[string]*Service),
	}
}

func (s *Services) New(baseURL, rootDir string, handlers map[string]Handler, handleDir Handler) *Service {
	baseURL = "/" + strings.Trim(baseURL, "/")
	rootDir = strings.TrimRight(rootDir, "/")
	srv := &Service{
		BaseURL:   baseURL,
		RootDir:   rootDir,
		parent:    s,
		Handlers:  handlers,
		HandleDir: handleDir,
	}
	s.byURL[baseURL] = srv
	return srv
}

func (s *Services) Run(addr string) {
	conn := service.MustClient(addr)

	s.Register(conn)

	conn.Run()
}

func (s *Services) Register(conn *service.Client) {
	for _, srv := range s.byURL {
		srv.Register(conn)
	}
}

func (s *Services) Handler(req *service.Request) *service.Response {
	var srv *Service
	sIdx := strings.IndexRune(req.Path[1:], '/')
	if sIdx == -1 {
		srv = s.byURL[req.Path]
	} else {
		srv = s.byURL[req.Path[:sIdx+1]]
	}
	if srv == nil {
		return req.ResponseErr(nil, http.StatusNotFound)
	}
	return srv.Handler(req)
}
