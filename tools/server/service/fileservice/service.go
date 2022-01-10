package fileservice

import (
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

	f, err := os.Open(path)
	if err != nil {
		return req.ResponseError(err, 404)
	}

	var h Handler = HandleBinary
	if info, err := f.Stat(); err != nil {
		return req.ResponseError(err, 500)
	} else if info.IsDir() {
		if s.HandleDir != nil {
			h = s.HandleDir
		} else {
			h = HandleDir
		}
	} else {
		ext := filepath.Ext(path)
		sh, found := s.Handlers[ext]
		if found {
			h = sh
		} else {
			ph, found := s.parent.Handlers[ext]
			if found {
				h = ph
			}
		}
	}
	return h(path, f, req)
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

func (s *Services) New(baseURL, rootDir string) *Service {
	baseURL = "/" + baseURL
	srv := &Service{
		BaseURL: baseURL,
		RootDir: rootDir,
		parent:  s,
	}
	s.byURL[baseURL] = srv
	return srv
}

func (s *Services) Run(addr string) {
	conn := service.MustClient(addr)

	for _, srv := range s.byURL {
		srv.Register(conn)
	}

	conn.Run()
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
		return req.ResponseError(nil, http.StatusNotFound)
	}
	return srv.Handler(req)
}
